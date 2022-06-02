package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Client struct {
	Host       string
	Port       string
	OAuthToken string
}

// GetTopics fetches all the topics from a specified kafka cluster
func (client *Client) GetTopics() (map[string]kafka.TopicMetadata, error) {
	// Create admin client
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": client.Host})
	if err != nil {
		return nil, err
	}

	// Set OAUTH token
	if client.OAuthToken != "" {
		err = adminClient.SetOAuthBearerToken(kafka.OAuthBearerToken{TokenValue: client.OAuthToken})
		if err != nil {
			return nil, err
		}
	}

	// Get metadata
	metadata, err := adminClient.GetMetadata(nil, true, 100)
	if err != nil {
		return nil, err
	}

	return metadata.Topics, nil
}
