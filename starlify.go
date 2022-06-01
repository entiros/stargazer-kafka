package stargazer_kafka

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

func (cfg *Config) ReadServices() error {
	var system StarlifySystem

	client := resty.New()

	// Fetch information about services from Starlify
	// This examples fetches a specific attribute from a specific pre-created Starlify service
	resp, _ := client.R().
		SetHeader("X-API-KEY", cfg.Starlify.ApiKey).
		Get(cfg.Starlify.ApiUrl + cfg.Starlify.SystemID)

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
		log.Info().Msgf("creating topic as a Starlify service: %s", topic)
		if err != nil {
			return "", err
		}
	}
	log.Info().Msg("created services in Starlify")
	return "services created", nil
}
