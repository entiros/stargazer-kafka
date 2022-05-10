package stargazer_kafka

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

func (cfg *Config) ReadServices() error {
	var system StarlifySystem

	client := resty.New()

	// Fetch information about services from Starlify
	// This examples fetches a specific attribute from a specific pre-created Starlify service
	resp, _ := client.R().
		SetHeader("X-API-KEY", "0fe184ed-b79a-4c5b-920b-732d0e98a148.ybZTnx0iX6mQlmFV8qIROsa9F2kAYyi9TVAisjKHpbfZBeKcL5u5XGIqBLlyooB3").
		Get("https://api.starlify.com/hypermedia/systems/336c9813-f271-4503-a14e-acc5698009cd")

	err := json.Unmarshal(resp.Body(), &system)
	if err != nil {
		fmt.Println(err)
	}

	// Check if the "Kafka topic" attribute has a value
	if system.Description != "" {
		// Creates the kafka.go topic
		cfg.CreateTopic("user")
	}
	return nil
}

type ServiceBody struct {
	Name string `json:"name,omitempty"`
}

func (cfg *Config) CreateService(topics []string) (string, error) {

	client := resty.New()

	for _, topic := range topics {
		_, err := client.R().
			SetHeader("X-API-KEY", cfg.Starlify.ApiKey).
			SetBody(&ServiceBody{Name: topic}).
			Post(cfg.Starlify.ApiUrl + cfg.Starlify.SystemID + "/services")
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("creating topic(s) in Starlify"), nil
	}
	return "info: finished creating Starlify services", nil
}
