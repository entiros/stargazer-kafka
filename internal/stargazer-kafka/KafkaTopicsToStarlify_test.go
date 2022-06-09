package stargazer_kafka

import (
	"context"
	kafka2 "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/entiros/stargazer-kafka/internal/kafka"
	"github.com/entiros/stargazer-kafka/internal/starlify"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"gopkg.in/h2non/gock.v1"
	"log"
	"strings"
	"testing"
	"time"
)

func createStarlifyClient() *starlify.Client {
	starlifyClient := &starlify.Client{
		BaseUrl:  "http://127.0.0.1:8080/hypermedia",
		ApiKey:   "api-key-123",
		AgentId:  "agent-id-123",
		SystemId: "system-id-123",
	}

	// Intercept Starlify client
	gock.InterceptClient(starlifyClient.GetRestyClient().GetClient())

	return starlifyClient
}

func createKafkaClient() *kafka.Client {
	return &kafka.Client{
		Host:       "127.0.0.1:9092",
		OAuthToken: "",
	}
}

func createKafkaTopic(kafkaClient *kafka.Client, topic string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	adminClient, err := kafkaClient.GetAdminClient()
	if err != nil {
		panic(err)
	}

	// Wait for operation to finish or at most 10s
	maxDur, err := time.ParseDuration("10s")
	if err != nil {
		panic(err)
	}

	_, err = adminClient.CreateTopics(
		ctx,
		[]kafka2.TopicSpecification{{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		}},
		kafka2.SetAdminOperationTimeout(maxDur))
	if err != nil {
		panic(err)
	}
}

func createGock() *gock.Request {
	return gock.New("http://127.0.0.1:8080/hypermedia")
}

// waitForKafkaContainer will wait for the GetTopic function to execute successfully
func waitForKafkaContainer(kafkaClient *kafka.Client) {
	for i := 0; i < 60*2; i++ {
		_, err := kafkaClient.GetTopics()
		if err == nil {
			log.Println("Kafka container up and running!")
			return
		}
		log.Println("Waiting for Kafka container...")
		time.Sleep(1 * time.Second)
	}
}

func destroyKafkaContainer(compose *testcontainers.LocalDockerCompose) {
	compose.Down()
	time.Sleep(1 * time.Second)
}

// initialize will start Kafka container and return a Kafka Client, Starlify client and the Kafka container compose
func initialize() (*kafka.Client, *starlify.Client, *testcontainers.LocalDockerCompose) {
	// Create and start Kafka container
	kafkaContainer := testcontainers.NewLocalDockerCompose(
		[]string{"configs/kafka/test-docker-compose.yml"},
		strings.ToLower(uuid.New().String()),
	)
	kafkaContainer.WithCommand([]string{"up", "-d"}).Invoke()

	// Create client and wait for Kafka container to start up
	kafkaClient := createKafkaClient()
	waitForKafkaContainer(kafkaClient)

	starlifyClient := createStarlifyClient()

	return kafkaClient, starlifyClient, kafkaContainer
}

func TestKafkaTopicsToStarlify_createKafkaTopicsToStarlify(t *testing.T) {
	kafkaClient, starlifyClient, kafkaContainer := initialize()
	defer destroyKafkaContainer(kafkaContainer)

	kafkaTopicsToStarlify := KafkaTopicsToStarlify{
		kafka:                   kafkaClient,
		starlify:                starlifyClient,
		topicsCache:             make(map[string]bool),
		lastUpdateReportedError: false,
	}

	createServiceRequest := createGock().
		Post("/systems/system-id-123/services")

	// For first execution, no topics is available
	kafkaTopicsToStarlify.createKafkaTopicsToStarlify()

	// No calls to POST services should have been made
	assert.True(t, gock.IsPending())
	gock.Remove(createServiceRequest.Mock)

	createKafkaTopic(kafkaClient, "Topic_1")
	createKafkaTopic(kafkaClient, "Topic_2")

	// Topic 1 is already existing in Starlify
	createGock().
		Get("/systems/system-id-123/services").
		Reply(200).
		JSON(starlify.ServicesPage{
			Services: []starlify.Service{
				{Id: "service-1", Name: "Topic_1"},
			},
			Page: starlify.Page{
				Size:          100,
				TotalElements: 1,
				TotalPages:    1,
				Number:        0,
			},
		})

	// Topic 2 will be created
	createGock().
		Post("/systems/system-id-123/services").
		MatchType("json").
		JSON(starlify.ServiceRequest{Name: "Topic_2"}).
		Reply(201).
		JSON(starlify.Service{Id: "service-1", Name: "Topic_2"})

	// Details will be updated
	createGock().
		Patch("/agents/agent-id-123").
		//MatchType("json").
		JSON(starlify.AgentRequest{Details: &starlify.Details{Topics: []starlify.TopicDetails{
			{
				Name: "Topic_1",
				Partitions: []starlify.PartitionDetails{
					{ID: 0},
				},
			},
			{
				Name: "Topic_2",
				Partitions: []starlify.PartitionDetails{
					{ID: 0},
				},
			},
		}}}).
		Reply(200).
		JSON(starlify.Agent{Id: "agent-id-123", Name: "Test agent", AgentType: "kafka"})

	// Topic 2 will be created
	kafkaTopicsToStarlify.createKafkaTopicsToStarlify()

	assert.True(t, gock.IsDone())
}
