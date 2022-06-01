package kafka

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
	"time"
)

type Client struct {
	Host string
	Port string
}

func (client *Client) CreateTopic(topic string) {
	a, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": client.Host})
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

// GetTopics fetches all the topics from a specified kafka cluster
func (client *Client) GetTopics() ([]string, error) {
	// Create admin client
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": client.Host})
	if err != nil {
		return nil, err
	}

	// Get metadata
	metadata, err := adminClient.GetMetadata(nil, true, 100)
	if err != nil {
		return nil, err
	}

	// Extract topic names
	var topics []string
	for _, topicMetadata := range metadata.Topics {
		topics = append(topics, topicMetadata.Topic)
	}

	return topics, nil
}
