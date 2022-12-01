package main

import (
	"log"
	"os"

	"github.com/entiros/stargazer-kafka/internal/kafka"
	stargazerkafka "github.com/entiros/stargazer-kafka/internal/stargazer-kafka"
	"github.com/entiros/stargazer-kafka/internal/starlify"
)

func main() {
	// Load configuration (from file or environment variables)
	config, err := stargazerkafka.LoadConfig()
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}

	// Starlify client
	starlifyClient := starlify.Client{
		BaseUrl:  config.Starlify.BaseUrl,
		ApiKey:   config.Starlify.ApiKey,
		AgentId:  config.Starlify.AgentId,
		SystemId: config.Starlify.SystemId,
	}

	// Kafka Client
	kafkaClient := kafka.Client{
		Host:       config.Kafka.Host,
		Type:       config.Kafka.Type,
		OAuthToken: config.Kafka.Auth.OAuth.Token,
	}

	// Create integration
	kafkaTopicsToStarlify, err := stargazerkafka.InitKafkaTopicsToStarlify(&kafkaClient, &starlifyClient)
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}

	// Start
	kafkaTopicsToStarlify.StartBlocking()
}
