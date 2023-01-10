package kafka

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/entiros/stargazer-kafka/internal/log"
	"github.com/twmb/franz-go/pkg/sasl"
	"github.com/twmb/franz-go/pkg/sasl/oauth"
	"github.com/twmb/franz-go/pkg/sasl/plain"
	"net"
	"time"

	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/aws/aws-sdk-go/aws/session"
	faws "github.com/twmb/franz-go/pkg/sasl/aws"
)

type Client struct {
	Hosts      []string
	AuthMethod sasl.Mechanism
	AWSSession func() *session.Session
}

func (k *Client) Client() (*kgo.Client, error) {
	return createClient(k.Hosts, k.AuthMethod)
}

func (k *Client) AdminClient() (*kadm.Client, error) {

	client, err := createClient(k.Hosts, k.AuthMethod)
	return kadm.NewClient(client), err

}

func NewKafkaClient(options ...func(*Client)) *Client {

	c := &Client{}
	for _, opt := range options {
		opt(c)
	}
	return c

}

func WithBootstrapServers(hosts ...string) func(*Client) {

	return func(client *Client) {
		client.Hosts = hosts
	}
}

func WithSession(foo func() *session.Session) func(client *Client) {
	return func(client *Client) {
		client.AWSSession = foo
	}
}

func Session() *session.Session {
	s, err := session.NewSession()
	if err != nil {
		log.Logger.Errorf("Failed to create session: %v", err)
	}
	return s
}

func WithOAuth(token string) func(client *Client) {
	return func(client *Client) {
		client.AuthMethod = oauth.Auth{
			Token: token,
		}.AsMechanism()
	}
}

func WithPassword(username string, password string) func(client *Client) {
	return func(client *Client) {
		client.AuthMethod = plain.Auth{
			User: username,
			Pass: password,
		}.AsMechanism()
	}
}

func WithIAM(accessKey string, secret string) func(client *Client) {

	return func(client *Client) {

		client.AuthMethod = faws.ManagedStreamingIAM(func(ctx context.Context) (faws.Auth, error) {
			sess := client.AWSSession()
			if sess == nil {
				return faws.Auth{}, fmt.Errorf("failed to create AWS session")
			}
			val, err := sess.Config.Credentials.GetWithContext(ctx)
			if err != nil {
				return faws.Auth{}, err
			}

			var accessKeyId string
			var secretAccessKey string
			if accessKey != "" {
				accessKeyId = accessKey
			} else {
				accessKeyId = val.AccessKeyID
			}

			if secret != "" {
				secretAccessKey = secret
			} else {
				secretAccessKey = val.SecretAccessKey
			}

			return faws.Auth{
				AccessKey:    accessKeyId,
				SecretKey:    secretAccessKey,
				SessionToken: val.SessionToken,
				UserAgent:    "entiros/kafka/v1.0.0",
			}, nil
		})
	}
}

func createClient(bootstrapServers []string, authMethod sasl.Mechanism) (*kgo.Client, error) {

	var opts []kgo.Opt

	opts = append(opts, kgo.SeedBrokers(bootstrapServers...))
	if authMethod != nil {
		opts = append(opts, kgo.SASL(authMethod))
		opts = append(opts, kgo.Dialer((&tls.Dialer{NetDialer: &net.Dialer{Timeout: 60 * time.Second}}).DialContext))
	}

	return kgo.NewClient(opts...)

}

func (c *Client) AllTopics(ctx context.Context) func() (kadm.TopicDetails, error) {

	return func() (kadm.TopicDetails, error) {
		return c.GetAllTopics(ctx)
	}

}

// GetTopics fetches all the topics from a specified kafka cluster.
func (c *Client) GetAllTopics(ctx context.Context) (kadm.TopicDetails, error) {

	// Get Kafka admin client
	client, err := c.AdminClient()
	if err != nil {
		return nil, err
	}

	// Get metadata
	metadata, err := client.Metadata(ctx)
	if err != nil {
		return nil, err
	}

	return metadata.Topics, nil
}

// GetTopics fetches all the topics from a specified kafka cluster.
func (c *Client) GetTopics(ctx context.Context) (kadm.TopicDetails, error) {

	// Get Kafka admin client
	client, err := c.AdminClient()
	if err != nil {
		return nil, err
	}

	// Get metadata
	metadata, err := client.Metadata(ctx)
	if err != nil {
		return nil, err
	}

	return metadata.Topics, nil
}

// GetTopics fetches all the topics from a specified kafka cluster.
func (c *Client) Ping(ctx context.Context) error {

	// Get Kafka admin client
	client, err := c.AdminClient()
	if err != nil {
		return nil
	}

	// Get metadata
	_, err = client.Metadata(ctx)

	return err
}

// CreateTopic fetches all the topics from a specified kafka cluster.
func (c *Client) CreateTopics(ctx context.Context, topics ...string) error {

	if len(topics) == 0 {
		return nil
	}
	// Get Kafka admin kafkaClient
	kafkaClient, err := c.AdminClient()
	if err != nil {
		return err
	}

	_, err = kafkaClient.CreateTopics(ctx, 1, 1, nil, topics...)
	if err != nil {
		return err
	}

	// TODO: Figure out what we should return here
	return nil
}

// CreateTopic fetches all the topics from a specified kafka cluster.
func (c *Client) DeleteTopics(ctx context.Context, topics ...string) error {

	if len(topics) == 0 {
		return nil
	}

	// Get Kafka admin kafkaClient
	kafkaClient, err := c.AdminClient()
	if err != nil {
		return err
	}

	_, err = kafkaClient.DeleteTopics(ctx, topics...)
	if err != nil {
		return err
	}

	// TODO: Figure out what we should return here
	return nil
}
