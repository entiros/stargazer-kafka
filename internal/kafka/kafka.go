package kafka

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/entiros/stargazer-kafka/internal/msk"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/aws/aws-sdk-go/aws/session"
	faws "github.com/twmb/franz-go/pkg/sasl/aws"
)

// Client holds the connection details for the Kafka instance.
type Client struct {
	Host        string
	Type        string
	OAuthToken  string
	adminClient *kadm.Client
}

func (c *Client) createAdminClient(bootstrapServers []string) (*kadm.Client, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	log.Printf("Getting admin client. Bootstrap: %v\n", bootstrapServers)
	kgoClient, err := kgo.NewClient(
		kgo.SeedBrokers(bootstrapServers...),

		kgo.SASL(faws.ManagedStreamingIAM(func(ctx context.Context) (faws.Auth, error) {
			val, err := sess.Config.Credentials.GetWithContext(ctx)
			if err != nil {
				return faws.Auth{}, err
			}

			return faws.Auth{
				AccessKey:    val.AccessKeyID,
				SecretKey:    val.SecretAccessKey,
				SessionToken: val.SessionToken,
				UserAgent:    "entiros/kafka/v1.0.0",
			}, nil
		})),

		kgo.Dialer((&tls.Dialer{NetDialer: &net.Dialer{Timeout: 10 * time.Second}}).DialContext),
	)
	if err != nil {
		return nil, err
	}

	return kadm.NewClient(kgoClient), nil
}

func (c *Client) createClient(bootstrapServers []string) (*kgo.Client, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	return kgo.NewClient(
		kgo.SeedBrokers(bootstrapServers...),

		kgo.SASL(faws.ManagedStreamingIAM(func(ctx context.Context) (faws.Auth, error) {
			val, err := sess.Config.Credentials.GetWithContext(ctx)
			if err != nil {
				return faws.Auth{}, err
			}

			return faws.Auth{
				AccessKey:    val.AccessKeyID,
				SecretKey:    val.SecretAccessKey,
				SessionToken: val.SessionToken,
				UserAgent:    "entiros/kafka/v1.0.0",
			}, nil
		})),

		kgo.Dialer((&tls.Dialer{NetDialer: &net.Dialer{Timeout: 10 * time.Second}}).DialContext),
	)
}

// GetTopics fetches all the topics from a specified kafka cluster.
func (c *Client) GetTopics() (kadm.TopicDetails, error) {
	servers, err := msk.GetMSKIamBootstrap(context.Background(), c.Host)
	if err != nil {
		fmt.Println("failed to get MSK bootstrap servers")
	}

	// Get Kafka admin client
	client, err := c.createAdminClient(servers)
	if err != nil {
		return nil, err
	}

	// Get metadata
	metadata, err := client.Metadata(nil)
	if err != nil {
		return nil, err
	}

	return metadata.Topics, nil
}

// CreateTopic fetches all the topics from a specified kafka cluster.
func (c *Client) CreateTopic(topicName string) (string, error) {
	ctx := context.Background()

	servers, err := msk.GetMSKIamBootstrap(ctx, c.Host)
	if err != nil {
		fmt.Println("failed to get MSK bootstrap servers")
	}

	// Get Kafka admin kafkaClient
	kafkaClient, err := c.createAdminClient(servers)
	if err != nil {
		return "", err
	}

	_, err = kafkaClient.CreateTopics(ctx, 1, 1, nil, topicName)
	if err != nil {
		return "", err
	}

	// TODO: Figure out what we should return here
	return "topic created", nil
}
