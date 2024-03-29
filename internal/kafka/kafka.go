package kafka

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/entiros/stargazer-kafka/internal/log"
	"github.com/entiros/stargazer-kafka/internal/prefix"
	"github.com/twmb/franz-go/pkg/sasl"
	"github.com/twmb/franz-go/pkg/sasl/oauth"
	"github.com/twmb/franz-go/pkg/sasl/plain"
	"net"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
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
		return nil
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

	os.Setenv("AWS_ACCESS_KEY", accessKey)
	os.Setenv("AWS_SECRET_KEY", secret)

	return func(client *Client) {

		client.AuthMethod = faws.ManagedStreamingIAM(func(ctx context.Context) (faws.Auth, error) {
			sess := client.AWSSession()
			if sess == nil {
				log.Logger.Errorf("Failed to create AWS session")
				return faws.Auth{}, fmt.Errorf("failed to create AWS session")
			}
			val, err := sess.Config.Credentials.GetWithContext(ctx)
			if err != nil {
				return faws.Auth{}, err
			}

			log.Logger.Debugf("Using Kafka IAM client %s", val.ProviderName)

			return faws.Auth{
				AccessKey:    val.AccessKeyID,
				SecretKey:    val.SecretAccessKey,
				SessionToken: val.SessionToken,
				UserAgent:    "entiros/kafka/v1.0.0",
			}, nil

		})
	}
}

func createClient(bootstrapServers []string, authMethod sasl.Mechanism) (*kgo.Client, error) {

	var opts []kgo.Opt

	opts = append(opts, kgo.SeedBrokers(bootstrapServers...))
	//	opts = append(opts, kgo.WithLogger(kzap.New(log.Logger.Desugar())))
	if authMethod != nil {
		opts = append(opts, kgo.SASL(authMethod))
		opts = append(opts, kgo.Dialer((&tls.Dialer{NetDialer: &net.Dialer{Timeout: 60 * time.Second}}).DialContext))
	}

	return kgo.NewClient(opts...)

}

func (c *Client) GetPrefixes(ctx context.Context) ([]string, error) {

	return Prefixes(ctx, c)

}

func Prefixes(ctx context.Context, c *Client) ([]string, error) {
	topics, err := c.GetTopics(ctx)
	if err != nil {
		return nil, err
	}

	var prefixMap map[string]bool
	for _, topic := range topics {
		pre, err := prefix.ExtractPrefix(topic.Topic)
		if err != nil {
			continue
		}
		prefixMap[pre] = true
	}

	var prefixes []string
	for prefix := range prefixMap {
		prefixes = append(prefixes, prefix)
	}

	return prefixes, nil
}

// GetTopics fetches all the topics from a specified kafka cluster.
func (c *Client) GetTopics(ctx context.Context) (kadm.TopicDetails, error) {

	return Topics(ctx, c)
}

func Topics(ctx context.Context, c *Client) (kadm.TopicDetails, error) {

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	client, err := c.AdminClient()
	if err != nil {
		log.Logger.Debugf("Failed to create Kafka Admin client. %v", err)
		return nil, err
	}
	defer client.Close()

	metadata, err := client.Metadata(ctx)
	if err != nil {
		log.Logger.Debugf("Failed to get metadata: %v", err)
		return nil, err
	}

	return metadata.Topics, nil
}

// CreateTopic fetches all the topics from a specified kafka cluster.
func (c *Client) CreateTopics(ctx context.Context, topics ...string) error {

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	if len(topics) == 0 {
		return nil
	}
	// Get Kafka admin kafkaClient
	kafkaClient, err := c.AdminClient()
	if err != nil {
		return err
	}
	defer kafkaClient.Close()

	_, err = kafkaClient.CreateTopics(ctx, 1, 1, nil, topics...)
	if err != nil {
		return err
	}

	// TODO: Figure out what we should return here
	return nil
}

// CreateTopic fetches all the topics from a specified kafka cluster.
func (c *Client) DeleteTopics(ctx context.Context, topics ...string) error {

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	if len(topics) == 0 {
		return nil
	}

	// Get Kafka admin kafkaClient
	kafkaClient, err := c.AdminClient()
	if err != nil {
		return err
	}
	defer kafkaClient.Close()

	_, err = kafkaClient.DeleteTopics(ctx, topics...)
	if err != nil {
		return err
	}

	// TODO: Figure out what we should return here
	return nil
}
