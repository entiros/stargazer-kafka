package kafka

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Client struct {
	Host        string
	OAuthToken  string
	adminClient *kafka.AdminClient
}

// getAdminClient will return Kafka Admin Client
func (c *Client) getAdminClient() (*kafka.AdminClient, error) {
	if c.adminClient == nil {

		// Create admin _this
		adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": c.Host})
		if err != nil {
			return nil, err
		}

		// Set OAUTH token
		if c.OAuthToken != "" {
			err = adminClient.SetOAuthBearerToken(kafka.OAuthBearerToken{TokenValue: c.OAuthToken})
			if err != nil {
				return nil, err
			}
		}

		c.adminClient = adminClient
	}

	return c.adminClient, nil
}

// GetTopics fetches all the topics from a specified kafka cluster
func (c *Client) GetTopics() (map[string]kafka.TopicMetadata, error) {
	// Get Kafka admin client
	client, err := c.getAdminClient()
	if err != nil {
		return nil, err
	}

	// Get metadata
	metadata, err := client.GetMetadata(nil, true, 100)
	if err != nil {
		return nil, err
	}

	return metadata.Topics, nil
}

// GetTopics fetches all the topics from a specified kafka cluster
func (c *Client) CreateTopic(topicName string) (string, error) {
	// Get Kafka admin client
	client, err := c.getAdminClient()
	if err != nil {
		return "", err
	}

	// Creatae a topic in Kafka
	_, err = client.CreateTopics(context.Background(), []kafka.TopicSpecification{
		{Topic: topicName, NumPartitions: 1},
	})
	if err != nil {
		return "", err
	}

	// TODO: Figure out what we should return here
	return "topic created", nil
}
