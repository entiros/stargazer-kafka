package main

import (
	"github.com/entiros/stargazer-kafka/internal/kafka"
	stargazerkafka "github.com/entiros/stargazer-kafka/internal/stargazer-kafka"
	"github.com/entiros/stargazer-kafka/internal/starlify"
	"log"
	"os"
)

func main() {
	//stargazerkafka.CreateTopic("test1")
	//stargazerkafka.CreateTopic("test2")
	//stargazerkafka.CreateTopic("test3")

	// Load configuration
	config, err := stargazerkafka.LoadConfig(os.Getenv("configFile"))
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}

	// Starlify client
	var starlifyClient = starlify.Client{
		BaseUrl:  config.Starlify.BaseUrl,
		ApiKey:   config.Starlify.ApiKey,
		AgentId:  config.Starlify.AgentId,
		SystemId: config.Starlify.SystemId,
	}

	// Kafka Client
	var kafkaClient = kafka.Client{
		Host: config.Kafka.Host,
		Port: config.Kafka.Port,
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
