package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Client struct {
	Host        string
	OAuthToken  string
	adminClient *kafka.AdminClient
}

// getAdminClient will return Kafka Admin Client
func (_this *Client) getAdminClient() (*kafka.AdminClient, error) {
	if _this.adminClient == nil {

		// Create admin _this
		adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": _this.Host})
		if err != nil {
			return nil, err
		}

		// Set OAUTH token
		if _this.OAuthToken != "" {
			err = adminClient.SetOAuthBearerToken(kafka.OAuthBearerToken{TokenValue: _this.OAuthToken})
			if err != nil {
				return nil, err
			}
		}

		_this.adminClient = adminClient
	}

	return _this.adminClient, nil
}

// GetTopics fetches all the topics from a specified kafka cluster
func (_this *Client) GetTopics() (map[string]kafka.TopicMetadata, error) {
	// Get Kafka admin client
	client, err := _this.getAdminClient()
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
