package stargazer_kafka

import (
	"github.com/entiros/stargazer-kafka/internal/config"
	"sort"
	"testing"

	"github.com/entiros/stargazer-kafka/internal/starlify"
	"gopkg.in/h2non/gock.v1"
)

func TestLoadConfig(t *testing.T) {

	c, err := config.LoadConfig("config.yml")
	if err != nil {
		t.Log(err)
	}
	t.Log(c)

}

func TestSyncList(t *testing.T) {

	var source = []string{"A", "G", "A", "A"}
	var target = []string{"A", "B", "E", "F"}

	sort.Strings(source)
	sort.Strings(target)

	t.Logf("Source: %v", source)
	t.Logf("Target : %v", target)

	insert, delete := ListDiff(source, target)

	t.Logf("Insert: %v", insert)
	t.Logf("Delete : %v", delete)

}

func createStarlifyClient() *starlify.Client {
	starlifyClient := &starlify.Client{
		BaseUrl:      "http://127.0.0.1:8080/hypermedia",
		ApiKey:       "api-key-123",
		AgentId:      "agent-id-123",
		MiddlewareId: "system-id-123",
	}

	// Intercept Starlify client
	gock.InterceptClient(starlifyClient.GetRestyClient().GetClient())

	return starlifyClient
}

func createGock() *gock.Request {
	return gock.New("http://127.0.0.1:8080/hypermedia")
}

/*
func TestKafkaTopicsToStarlify_createKafkaTopicsToStarlify(t *testing.T) {
	defer gock.Off()

	starlifyClient := createStarlifyClient()

	kafkaTopicsToStarlify := KafkaTopicsToStarlify{
		starlify:                starlifyClient,
		topicsCache:             make(map[string]bool),
		lastUpdateReportedError: false,
	}

	createServiceRequest := createGock().
		Post("/systems/system-id-123/services")

	// For first execution, no topics is available
	kafkaTopicsToStarlify.createTopics(func() (map[string]kafka.TopicMetadata, error) {
		return map[string]kafka.TopicMetadata{}, nil
	})

	// No calls to POST services should have been made
	assert.True(t, gock.IsPending())
	gock.Remove(createServiceRequest.Mock)

	// 'Topic 1' is already existing in Starlify
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

	// 'Topic 2' will be created
	createGock().
		Post("/systems/system-id-123/services").
		MatchType("json").
		JSON(starlify.ServiceRequest{Name: "Topic_2"}).
		Reply(201).
		JSON(starlify.Service{Id: "service-1", Name: "Topic_2"})

	// Details will be updated
	createGock().
		Patch("/agents/agent-id-123").
		MatchType("json").
		JSON(starlify.AgentRequest{
			Error: "",
			Details: &starlify.Details{
				Topics: []starlify.TopicDetails{
					{
						Name:       "Topic_1",
						Partitions: []starlify.PartitionDetails{{ID: 0}},
					},
					{
						Name:       "Topic_2",
						Partitions: []starlify.PartitionDetails{{ID: 0}},
					},
				},
			},
		}).
		Reply(200).
		JSON(starlify.Agent{Id: "agent-id-123", Name: "Test agent", AgentType: "kafka"})

	// Run with two topics available in Kafka
	kafkaTopicsToStarlify.createTopics(func() (map[string]kafka.TopicMetadata, error) {
		return map[string]kafka.TopicMetadata{
			"Topic_1": {
				Topic: "Topic_1",
				Partitions: []kafka.PartitionMetadata{
					{
						ID: 0,
					},
				},
			},
			"Topic_2": {
				Topic: "Topic_2",
				Partitions: []kafka.PartitionMetadata{
					{
						ID: 0,
					},
				},
			},
		}, nil
	})

	assert.True(t, gock.IsDone())
}
*/
