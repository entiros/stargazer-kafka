package stargazer_kafka

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"os"
	"time"
)

func (cfg *Config) CreateTopic(topic string) {
	a, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": cfg.Kafka.Host})
	if err != nil {
		fmt.Printf("Failed to create Admin client: %s\n", err)
		os.Exit(1)
	}

	// Contexts are used to abort or limit the amount of time
	// the Admin call blocks waiting for a result.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create topics on cluster.
	// Set Admin options to wait for the operation to finish (or at most 60s)
	maxDur, err := time.ParseDuration("60s")
	if err != nil {
		panic("ParseDuration(60s)")
	}
	_, err = a.CreateTopics(
		ctx,
		// Multiple topics can be created simultaneously
		// by providing more TopicSpecification structs here.
		[]kafka.TopicSpecification{{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1}},
		// Admin options
		kafka.SetAdminOperationTimeout(maxDur))
	if err != nil {
		fmt.Printf("Failed to create topic: %v\n", err)
		os.Exit(1)
	}
}

// ReadTopics fetches all the topics from a specified kafka cluster
func (cfg *Config) ReadTopics() {
	var topics []string
	k, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": cfg.Kafka.Host})
	if err != nil {
		log.Fatal("failed to connect to kafka cluster")
	}

	metadata, err := k.GetMetadata(nil, true, 100)
	if err != nil {
		log.Print("failed to get metadata from kafka cluster")
	}

	for _, s := range metadata.Topics {
		topics = append(topics, s.Topic)
	}

	_, err = cfg.CreateService(topics)
	if err != nil {
		fmt.Printf("failed to create starlify service: %s\n", err)
		return
	}
}
