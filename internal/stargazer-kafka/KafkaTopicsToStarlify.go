package stargazer_kafka

import (
	"context"
	"fmt"
	"github.com/twmb/franz-go/pkg/kadm"
	"log"
	"sort"
	"strings"

	"github.com/entiros/stargazer-kafka/internal/kafka"
	"github.com/entiros/stargazer-kafka/internal/starlify"
)

type KafkaTopicsToStarlify struct {
	starlify                *starlify.Client
	kafka                   *kafka.Client
	topicsCache             map[string]bool
	lastUpdateReportedError bool
}

func InitKafkaTopicsToStarlify(ctx context.Context, kafkaClient *kafka.Client, starlify *starlify.Client) (*KafkaTopicsToStarlify, error) {
	// Get agent from Starlify and verify type
	agent, err := starlify.GetAgent(ctx)
	if err != nil {
		return nil, err
	} else if agent.AgentType != "managed-kafka" {
		return nil, fmt.Errorf("starlify agent '%s' is of type '%s', expected 'kafka'", agent.Id, agent.AgentType)
	}

	/*
	   TODO: what to replace this with?
	   	// Get services from Starlify to verify that system exists
	   	_, err = starlify.GetServices(ctx)
	   	if err != nil {
	   		return nil, err
	   	}
	*/
	// Get topics from Kafka to verify connection
	err = kafkaClient.Ping(ctx)
	if err != nil {
		return nil, err
	}

	kafkaTopicsToStarlify := KafkaTopicsToStarlify{
		starlify:                starlify,
		kafka:                   kafkaClient,
		topicsCache:             make(map[string]bool),
		lastUpdateReportedError: false,
	}

	return &kafkaTopicsToStarlify, nil
}

// reportError will send error message to Starlify and mark last update to have reported error.
func (k *KafkaTopicsToStarlify) ReportError(ctx context.Context, message error) {
	log.Println(message)

	// // Report error
	err := k.starlify.ReportError(ctx, message.Error())
	if err != nil {
		log.Printf("Failed to report error: %v", err)
	} else {
		k.lastUpdateReportedError = true
	}
}

func (k *KafkaTopicsToStarlify) Ping(ctx context.Context) error {
	err := k.starlify.Ping(ctx)
	if err != nil {
		log.Printf("Starlify Ping failed: %s", err.Error())
	}
	return err
}

// We are not supposed to understand this.
// Read here:
//  http://www.mlsite.net/blog/?p=2250
//  http://www.mlsite.net/blog/?p=2271

type tuple struct {
	position int
	value    string
}

// Produce two lists required to make source equal to target.
// - inserts, what need to be inserted into source
// - deletes, what need to be deleted from source
func ListDiff(source []string, target []string) ([]string, []string) {

	sort.Strings(source)
	sort.Strings(target)

	ins, del := SyncActions(source, target)

	var inserts []string
	for _, r := range ins {
		inserts = append(inserts, r.value)
	}

	var deletes []string
	for _, d := range del {
		deletes = append(deletes, source[d])
	}

	return inserts, deletes
}

func SyncActions(a []string, target []string) ([]tuple, []int) {

	var inserts []tuple
	var deletes []int
	x := 0
	y := 0

	for (x < len(a)) || (y < len(target)) {
		if y >= len(target) {
			deletes = append(deletes, x)
			x += 1
		} else if x >= len(a) {
			inserts = append(inserts, tuple{y, target[y]})
			y += 1
		} else if a[x] < target[y] {
			deletes = append(deletes, x)
			x += 1
		} else if a[x] > target[y] {
			inserts = append(inserts, tuple{y, target[y]})
			y += 1
		} else {
			x += 1
			y += 1
		}
	}
	return inserts, deletes
}

func (k *KafkaTopicsToStarlify) getStarlifyTopics(ctx context.Context) (string, []string, error) {

	log.Printf("Getting Starlify topics")
	prefix, topics, err := k.starlify.GetTopics(ctx)
	if err != nil {
		return "", nil, err
	}

	log.Printf("%d topics receives from Starlify: %v", len(topics), topics)
	return prefix, topics, nil
}

// createTopics will get topics from Kafka and create matching services in Starlify.
func (k *KafkaTopicsToStarlify) SyncTopics(ctx context.Context) error {

	prefix, starlifyTopics, err := k.getStarlifyTopics(ctx)
	if err != nil {
		return err
	}

	kafkaTopics, err := k.getKafkaTopics(ctx, prefix)
	if err != nil {
		return err
	}

	createMe, deleteMe := ListDiff(kafkaTopics, starlifyTopics)

	log.Printf("Creating topics: %v", createMe)
	err = k.kafka.CreateTopics(ctx, createMe...)
	if err != nil {
		return err
	}

	log.Printf("Deleting topics: %v", deleteMe)
	err = k.kafka.DeleteTopics(ctx, deleteMe...)
	if err != nil {
		return err
	}

	/*
		// Update details
		err = k.starlify.UpdateDetails(starlify.Details{Topics: createTopicDetails(topics)})
		if err != nil {
			e := fmt.Errorf("Failed to update details with error: %v", err.Error())
		}
	*/

	/*
		// Clear any errors reported
		if k.lastUpdateReportedError {
			err = k.starlify.ClearError()
			if err != nil {
				log.Printf("Failed to clear error")
			} else {
				k.lastUpdateReportedError = false
			}
		}
	*/

	return nil
}

func (k *KafkaTopicsToStarlify) getKafkaTopics(ctx context.Context, prefix string) ([]string, error) {

	log.Println("Fetching topics from Kafka")
	topics, err := k.kafka.GetTopics(ctx)
	if err != nil {
		err = fmt.Errorf("Failed to get topics from Kafka with error: %v" + err.Error())
		return nil, err
	}

	var kafkaTopics []string
	for _, t := range topics {
		if strings.HasPrefix(t.Topic, prefix) {
			kafkaTopics = append(kafkaTopics, t.Topic)
		}
	}
	log.Printf("%d topics received from Kafka: %v", len(kafkaTopics), kafkaTopics)

	return kafkaTopics, nil
}

// createTopicDetails will return topic details in Starlify format.
func createTopicDetails(topics kadm.TopicDetails) []starlify.TopicDetails {

	var topicDetails []starlify.TopicDetails

	topics.EachPartition(func(detail kadm.PartitionDetail) {
		topicDetails = append(topicDetails, starlify.TopicDetails{
			Name: detail.Topic,
			Partitions: []starlify.PartitionDetails{
				{ID: detail.Partition},
			},
		})
	})

	// Sort by name
	sort.Slice(topicDetails, func(i, j int) bool {
		return topicDetails[i].Name < topicDetails[j].Name
	})

	return topicDetails
}
