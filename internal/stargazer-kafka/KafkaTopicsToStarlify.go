package stargazer_kafka

import (
	"fmt"
	ck "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/entiros/stargazer-kafka/internal/kafka"
	"github.com/entiros/stargazer-kafka/internal/starlify"
	"github.com/go-co-op/gocron"
	"log"
	"time"
)

type KafkaTopicsToStarlify struct {
	kafkaClient             *kafka.Client
	starlifyClient          *starlify.Client
	kafkaTopicsCache        map[string]bool
	lastUpdateReportedError bool
}

func InitKafkaTopicsToStarlify(kafkaClient *kafka.Client, starlifyClient *starlify.Client) (*gocron.Scheduler, error) {
	agent, err := starlifyClient.GetAgent()
	if err != nil {
		return nil, err
	} else if agent.AgentType != "kafka" {
		return nil, fmt.Errorf("starlify agent '%s' is of type '%s', expected 'kafka'", agent.Id, agent.AgentType)
	}

	var kafkaTopicsToStarlify = KafkaTopicsToStarlify{
		kafkaClient:             kafkaClient,
		starlifyClient:          starlifyClient,
		kafkaTopicsCache:        make(map[string]bool),
		lastUpdateReportedError: false,
	}

	cron := gocron.NewScheduler(time.UTC)

	// Ping every 5 min
	_, err = cron.
		Every(5).
		Minutes().
		Do(starlifyClient.Ping)
	if err != nil {
		return nil, err
	}

	// Create Kafka topics in Starlify as services every minute
	_, err = cron.
		Every(1).
		Minutes().
		Do(kafkaTopicsToStarlify.createKafkaTopicsToStarlify)
	if err != nil {
		return nil, err
	}

	return cron, nil
}

// reportError will send error message to Starlify and mark last update to have reported error
func (_this *KafkaTopicsToStarlify) reportError(message string) {
	log.Println(message)

	// // Report error
	err := _this.starlifyClient.ReportError(message)
	if err != nil {
		log.Printf("Failed to report error")
	} else {
		_this.lastUpdateReportedError = true
	}
}

// createKafkaTopicsToStarlify will get topics from Kafka and create matching services in Starlify
func (_this *KafkaTopicsToStarlify) createKafkaTopicsToStarlify() {
	log.Println("Fetching topics from Kafka")

	topics, err := _this.kafkaClient.GetTopics()
	if err != nil {
		_this.reportError("Failed to get topics from Kafka")
		return
	}

	// Get topics not in cache
	var topicsToBeCreated []string
	for topic := range topics {
		if !_this.kafkaTopicsCache[topic] {
			topicsToBeCreated = append(topicsToBeCreated, topic)
		}
	}

	// There are possible topics to be created
	if len(topicsToBeCreated) > 0 {
		// Get Starlify services
		services, err := _this.starlifyClient.GetServices()
		if err != nil {
			_this.reportError("Failed to get Starlify services")
			return
		}

		// Create map of services in Starlify
		starlifyServices := make(map[string]bool)
		for _, service := range services {
			starlifyServices[service.Name] = true
		}

		for _, topic := range topicsToBeCreated {
			// Create service if it doesn't exist in Starlify
			if !starlifyServices[topic] {
				_, err := _this.starlifyClient.CreateService(topic)
				if err != nil {
					log.Printf("Failed to create service '%s'", topic)
					continue
				}
			}

			// Add to cache
			_this.kafkaTopicsCache[topic] = true
		}

		// Update details
		err = _this.starlifyClient.UpdateDetails(starlify.Details{Topics: createTopicDetails(topics)})
		if err != nil {
			_this.reportError("Failed to update details")
			return
		}
	} else {
		log.Println("No new service topics to create in Starlify")
	}

	// Clear any errors reported
	if _this.lastUpdateReportedError {
		err = _this.starlifyClient.ClearError()
		if err != nil {
			log.Printf("Failed to clear error")
		} else {
			_this.lastUpdateReportedError = false
		}
	}
}

func createTopicDetails(topics map[string]ck.TopicMetadata) []starlify.DetailsTopic {
	var topicDetails = make([]starlify.DetailsTopic, len(topics))

	var i = 0
	for _, topic := range topics {
		topicDetails[i] = starlify.DetailsTopic{
			Name:       topic.Topic,
			Partitions: createPartitionDetails(topic.Partitions),
		}
		i++
	}

	return topicDetails
}

func createPartitionDetails(partitions []ck.PartitionMetadata) []starlify.DetailsPartition {
	var partitionDetails = make([]starlify.DetailsPartition, len(partitions))

	var i = 0
	for _, partition := range partitions {
		partitionDetails[i] = starlify.DetailsPartition{
			ID: partition.ID,
		}
		i++
	}

	return partitionDetails
}
