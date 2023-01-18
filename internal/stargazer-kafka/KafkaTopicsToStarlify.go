package stargazer_kafka

import (
	"context"
	"fmt"
	"github.com/entiros/stargazer-kafka/internal/kafka"
	"github.com/entiros/stargazer-kafka/internal/log"
	pre "github.com/entiros/stargazer-kafka/internal/prefix"
	"github.com/entiros/stargazer-kafka/internal/starlify"
	"github.com/twmb/franz-go/pkg/kadm"
	"sort"
	"strings"
)

type KafkaTopicsToStarlify struct {
	starlify                *starlify.Client
	kafka                   *kafka.Client
	lastUpdateReportedError bool
}

const KafkaType = "managed-kafka"

func InitKafkaTopicsToStarlify(ctx context.Context, kafkaClient *kafka.Client, starlify *starlify.Client) (*KafkaTopicsToStarlify, error) {
	// Get agent from Starlify and verify type
	agent, err := starlify.GetAgent(ctx)
	if err != nil {
		return nil, err
	} else if agent.AgentType != KafkaType {
		return nil, fmt.Errorf("starlify agent '%s' is of type '%s', expected '%s'", agent.Id, agent.AgentType, KafkaType)
	}

	kafkaTopicsToStarlify := KafkaTopicsToStarlify{
		starlify:                starlify,
		kafka:                   kafkaClient,
		lastUpdateReportedError: false,
	}

	return &kafkaTopicsToStarlify, nil
}

// reportError will send error message to Starlify and mark last update to have reported error.
func (k *KafkaTopicsToStarlify) ReportError(ctx context.Context, message error) {
	log.Logger.Debug(message)

	// // Report error
	err := k.starlify.ReportError(ctx, message.Error())
	if err != nil {
		log.Logger.Debugf("Failed to report error: %v", err)
	} else {
		k.lastUpdateReportedError = true
	}
}

func (k *KafkaTopicsToStarlify) Ping(ctx context.Context) error {
	err := k.starlify.Ping(ctx)
	if err != nil {
		log.Logger.Debugf("Starlify Ping failed: %s", err.Error())
	}
	return err
}

func (k *KafkaTopicsToStarlify) getStarlifyTopics(ctx context.Context) (string, []starlify.TopicEndpoint, error) {

	log.Logger.Debugf("Getting Starlify topics")
	return k.starlify.GetTopics(ctx)

}

// get topics(endpoints on a middleware) from Starlify and create matching topics in Kafka.
func (k *KafkaTopicsToStarlify) SyncTopicsToKafka(ctx context.Context) (string, error) {

	// Get topics (endpoints) for this specific Middleware
	prefix, topics, err := k.getStarlifyTopics(ctx)
	if err != nil {
		return "", err
	}

	var starlifyTopics []string
	for _, topic := range topics {
		starlifyTopics = append(starlifyTopics, topic.Name)
	}

	log.Logger.Debugf("Prefix is: %s", prefix)
	if prefix == "" || len(prefix) < 8 {
		return "", fmt.Errorf("invalid prefix: %s", prefix)
	}

	// Get all Kafka topics with the specified prefix. Prefix is from Starlify middleware.
	kafkaTopics, err := k.getKafkaTopics(ctx, prefix)
	if err != nil {
		return "", err
	}

	createMe, deleteMe := ListDiff(kafkaTopics, starlifyTopics)

	log.Logger.Debugf("Creating topics: %v", createMe)
	err = k.kafka.CreateTopics(ctx, createMe...)
	if err != nil {
		return "", err
	}

	log.Logger.Debugf("Deleting topics: %v", deleteMe)
	err = k.kafka.DeleteTopics(ctx, deleteMe...)
	if err != nil {
		return "", err
	}

	return prefix, nil
}

// get topics from Kafka and create matching topics in Starlify.
func (k *KafkaTopicsToStarlify) SyncTopicsToStarlify(ctx context.Context) (string, error) {

	prefix, topics, err := k.getStarlifyTopics(ctx)
	if err != nil {
		return "", err
	}

	var starlifyTopics []string
	topicEndpoints := make(map[string]starlify.TopicEndpoint)

	for _, topic := range topics {
		starlifyTopics = append(starlifyTopics, topic.Name)
		topicEndpoints[topic.Name] = topic
	}

	kafkaTopics, err := k.getKafkaTopics(ctx, prefix)
	if err != nil {
		return "", err
	}

	createMe, deleteMe := ListDiff(starlifyTopics, kafkaTopics)

	log.Logger.Debugf("Creating topics: %v", createMe)
	for _, topic := range createMe {
		err = k.starlify.CreateTopic(ctx, topic)
		if err != nil {
			return "", err
		}
	}

	log.Logger.Debugf("Deleting topics: %v", deleteMe)
	for _, topic := range deleteMe {
		err = k.starlify.DeleteTopic(ctx, topicEndpoints[topic])
		if err != nil {
			return "", err
		}
	}
	return prefix, nil
}

func (k *KafkaTopicsToStarlify) getKafkaTopics(ctx context.Context, prefix string) ([]string, error) {

	if err := pre.Validate(prefix); err != nil {
		return nil, err
	}

	topics, err := k.kafka.GetTopics(ctx)
	if err != nil {
		err = fmt.Errorf("failed to get topics from Kafka with error: %v", err.Error())
		return nil, err
	}

	var kafkaTopics []string
	for _, t := range topics {
		if strings.HasPrefix(t.Topic, prefix) {
			kafkaTopics = append(kafkaTopics, t.Topic)
		}
	}
	log.Logger.Debugf("%d topics received from Kafka: %v", len(kafkaTopics), kafkaTopics)

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
