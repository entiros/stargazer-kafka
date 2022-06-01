package stargazer_kafka

import (
	"github.com/entiros/stargazer-kafka/internal/kafka"
	"github.com/entiros/stargazer-kafka/internal/starlify"
	"github.com/go-co-op/gocron"
	"log"
	"time"
)

type KafkaTopicsToStarlify struct {
	kafkaTopicsCache map[string]bool
}

func InitKafkaTopicsToStarlify(kafka *kafka.Client, starlify *starlify.Client) (*gocron.Scheduler, error) {
	var kafkaTopicsToStarlify = KafkaTopicsToStarlify{}
	kafkaTopicsToStarlify.kafkaTopicsCache = make(map[string]bool)

	//kafka.CreateTopic("Topic1")
	//kafka.CreateTopic("Topic2")
	//kafka.CreateTopic("Topic3")

	cron := gocron.NewScheduler(time.UTC)

	// Ping every 5 min
	_, err := cron.
		Every(5).
		Minutes().
		Do(starlify.Ping)
	if err != nil {
		return nil, err
	}

	// Create Kafka topics in Starlify as services every minute
	_, err = cron.
		Every(1).
		Minutes().
		Do(kafkaTopicsToStarlify.createKafkaTopicsToStarlify, kafka, starlify)
	if err != nil {
		return nil, err
	}

	return cron, nil
}

// createKafkaTopicsToStarlify will get topics from Kafka and create matching services in Starlify
func (_this *KafkaTopicsToStarlify) createKafkaTopicsToStarlify(kafka *kafka.Client, starlify *starlify.Client) error {
	topics, err := kafka.GetTopics()
	if err != nil {
		log.Println("Failed to get topics from Kafka")
		return err
	}

	// Get topics not in cache
	var topicsToBeCreated []string
	for _, topic := range topics {
		if !_this.kafkaTopicsCache[topic] {
			topicsToBeCreated = append(topicsToBeCreated, topic)
		}
	}

	// There are possible topics to be created
	if len(topicsToBeCreated) > 0 {
		// Get Starlify services
		services, err := starlify.GetServices()
		if err != nil {
			log.Println("Failed to get Starlify services")
			return err
		}

		// Create map of services in Starlify
		starlifyServices := make(map[string]bool)
		for _, service := range services {
			starlifyServices[service.Name] = true
		}

		for _, topic := range topicsToBeCreated {
			// Create service if it doesn't exist in Starlify
			if !starlifyServices[topic] {
				service, err := starlify.CreateService(topic)
				if err != nil {
					log.Printf("Failed to create service '%s'", topic)
				} else {
					_this.kafkaTopicsCache[service.Name] = true
				}
			}
		}
	}

	return nil
}
