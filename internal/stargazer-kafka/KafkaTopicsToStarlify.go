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
	starlify                *starlify.Client
	topicsCache             map[string]bool
	lastUpdateReportedError bool
}

func InitKafkaTopicsToStarlify(kafka *kafka.Client, starlify *starlify.Client) (*gocron.Scheduler, error) {
	// Get agent from Starlify and verify type
	agent, err := starlify.GetAgent()
	if err != nil {
		return nil, err
	} else if agent.AgentType != "kafka" {
		return nil, fmt.Errorf("starlify agent '%s' is of type '%s', expected 'kafka'", agent.Id, agent.AgentType)
	}

	// Get services from Starlify to verify that system exists
	_, err = starlify.GetServices()
	if err != nil {
		return nil, err
	}

	// Get topics from Kafka to verify connection
	_, err = kafka.GetTopics()
	if err != nil {
		return nil, err
	}

	var kafkaTopicsToStarlify = KafkaTopicsToStarlify{
		starlify:                starlify,
		topicsCache:             make(map[string]bool),
		lastUpdateReportedError: false,
	}

	cron := gocron.NewScheduler(time.UTC)

	// Ping to trigger last-seen to be updated
	_, err = cron.
		Every(5).
		Minutes().
		Do(kafkaTopicsToStarlify.ping)
	if err != nil {
		return nil, err
	}

	// Create Kafka topics in Starlify as services
	_, err = cron.
		Every(30).
		Seconds().
		Do(kafkaTopicsToStarlify.createKafkaTopicsToStarlify, kafka.GetTopics)
	if err != nil {
		return nil, err
	}

	return cron, nil
}

// reportError will send error message to Starlify and mark last update to have reported error
func (_this *KafkaTopicsToStarlify) reportError(message string) {
	log.Println(message)

	// // Report error
	err := _this.starlify.ReportError(message)
	if err != nil {
		log.Printf("Failed to report error")
	} else {
		_this.lastUpdateReportedError = true
	}
}

func (_this *KafkaTopicsToStarlify) ping() {
	err := _this.starlify.Ping()
	if err != nil {
		log.Printf("Starlify Ping failed: %s", err.Error())
	}
}

// createKafkaTopicsToStarlify will get topics from Kafka and create matching services in Starlify
func (_this *KafkaTopicsToStarlify) createKafkaTopicsToStarlify(getTopics func() (map[string]ck.TopicMetadata, error)) {
	log.Println("Fetching topics from Kafka")

	topics, err := getTopics()
	if err != nil {
		_this.reportError("Failed to get topics from Kafka with error: " + err.Error())
		return
	}

	log.Printf("%d topics received from Kafka", len(topics))

	// Get topics not in cache
	var topicsToBeCreated []string
	for topic := range topics {
		if !_this.topicsCache[topic] {
			topicsToBeCreated = append(topicsToBeCreated, topic)
		}
	}

	// There are possible topics to be created
	if len(topicsToBeCreated) > 0 {
		log.Printf("%d topics not in local cache and will be processed", len(topicsToBeCreated))

		// Get Starlify services
		services, err := _this.starlify.GetServices()
		if err != nil {
			_this.reportError("Failed to get Starlify services with error: " + err.Error())
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
				_, err := _this.starlify.CreateService(topic)
				if err != nil {
					log.Printf("Failed to create service '%s'", topic)
					continue
				}
			} else {
				log.Printf("Topic (service) '%s' already exists in Starlify", topic)
			}

			// Add to cache
			_this.topicsCache[topic] = true
		}

		// Update details
		err = _this.starlify.UpdateDetails(starlify.Details{Topics: createTopicDetails(topics)})
		if err != nil {
			_this.reportError("Failed to update details with error: " + err.Error())
			return
		}
	} else {
		log.Println("No new topics (services) to create in Starlify")
	}

	// Clear any errors reported
	if _this.lastUpdateReportedError {
		err = _this.starlify.ClearError()
		if err != nil {
			log.Printf("Failed to clear error")
		} else {
			_this.lastUpdateReportedError = false
		}
	}
}

// createTopicDetails will return topic details in Starlify format
func createTopicDetails(topics map[string]ck.TopicMetadata) []starlify.TopicDetails {
	var topicDetails = make([]starlify.TopicDetails, len(topics))

	var i = 0
	for _, topic := range topics {
		topicDetails[i] = starlify.TopicDetails{
			Name:       topic.Topic,
			Partitions: createPartitionDetails(topic.Partitions),
		}
		i++
	}

	return topicDetails
}

// createPartitionDetails will return topic partition details in Starlify format
func createPartitionDetails(partitions []ck.PartitionMetadata) []starlify.PartitionDetails {
	var partitionDetails = make([]starlify.PartitionDetails, len(partitions))

	var i = 0
	for _, partition := range partitions {
		partitionDetails[i] = starlify.PartitionDetails{
			ID: partition.ID,
		}
		i++
	}

	return partitionDetails
}
