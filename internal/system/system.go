package system

import (
	"context"
	"fmt"
	"github.com/entiros/stargazer-kafka/internal/config"
	"github.com/entiros/stargazer-kafka/internal/kafka"
	stargazerkafka "github.com/entiros/stargazer-kafka/internal/stargazer-kafka"
	"github.com/entiros/stargazer-kafka/internal/starlify"
	"log"
	"os"
	"time"
)

type System struct {
	cfg    *config.Config
	cancel context.CancelFunc
	ctx    context.Context
	file   string
}

var systems map[string]*System

func init() {
	systems = make(map[string]*System)
}

func AddSystem(ctx context.Context, c string) (*System, error) {

	if _, ok := systems[c]; !ok {
		return nil, fmt.Errorf("system %s exists already", c)
	}

	s, err := NewSystem(ctx, c)
	if err != nil {
		return nil, err
	}

	systems[c] = s
	go s.start(ctx)
	return s, nil
}

func DeleteSystem(ctx context.Context, c string) error {

	system, ok := systems[c]
	if !ok {
		return fmt.Errorf("could not find system for %s", c)
	}

	system.stop()
	delete(systems, c)

	return nil

}

func NewSystem(ctx context.Context, c string) (*System, error) {

	ctx, cancel := context.WithCancel(ctx)

	cfg, err := config.LoadConfig(c)
	if err != nil {
		return nil, err
	}

	s := &System{
		cfg:    cfg,
		cancel: cancel,
		ctx:    ctx,
		file:   c,
	}

	return s, nil

}

func (s *System) stop() {
	s.cancel()
}

func (s *System) start(ctx context.Context) {

	// Starlify client
	starlifyClient := starlify.Client{
		BaseUrl:      s.cfg.Starlify.BaseUrl,
		ApiKey:       s.cfg.Starlify.ApiKey,
		AgentId:      s.cfg.Starlify.AgentId,
		MiddlewareId: s.cfg.Starlify.MiddlewareId,
	}

	var kafkaClient *kafka.Client
	if s.cfg.Kafka.Auth.IAM.Secret != "" {
		kafkaClient = kafka.NewKafkaClient(
			kafka.WithBootstrapServers(s.cfg.Kafka.BootstrapServers...),
			kafka.WithSession(kafka.Session),
			kafka.WithIAM(s.cfg.Kafka.Auth.IAM.Key, s.cfg.Kafka.Auth.IAM.Secret),
		)
	} else if s.cfg.Kafka.Auth.Plain.User != "" {
		kafkaClient = kafka.NewKafkaClient(
			kafka.WithBootstrapServers(s.cfg.Kafka.BootstrapServers...),
			kafka.WithSession(kafka.Session),
			kafka.WithPassword(s.cfg.Kafka.Auth.Plain.User, s.cfg.Kafka.Auth.Plain.Password),
		)
	} else if s.cfg.Kafka.Auth.OAuth.Token != "" {
		kafkaClient = kafka.NewKafkaClient(
			kafka.WithBootstrapServers(s.cfg.Kafka.BootstrapServers...),
			kafka.WithSession(kafka.Session),
			kafka.WithOAuth(s.cfg.Kafka.Auth.OAuth.Token),
		)
	} else {
		kafkaClient = kafka.NewKafkaClient(
			kafka.WithBootstrapServers(s.cfg.Kafka.BootstrapServers...),
			kafka.WithSession(kafka.Session),
		)
	}

	// Create integration
	kafkaTopicsToStarlify, err := stargazerkafka.InitKafkaTopicsToStarlify(ctx, kafkaClient, &starlifyClient)
	if err != nil {
		log.Println(err)
		os.Exit(3)
	}

	doSync := time.After(time.Second * 5)
	doPing := time.After(time.Second * 1)

loop:
	for {
		select {

		case <-doPing:

			log.Println("Pinging Starlify")
			kafkaTopicsToStarlify.Ping(ctx)
			doPing = time.After(time.Second * 5)

		case <-doSync:

			log.Println("Performing Starlify-->Kafka topic sync")
			err := kafkaTopicsToStarlify.SyncTopics(ctx)
			if err != nil {
				log.Printf("Failed to sync topic. Retrying in 5m. %v", err)
				kafkaTopicsToStarlify.ReportError(ctx, err)
			}
			doSync = time.After(time.Second * 20)

		case <-ctx.Done():
			break loop
		}
	}

	log.Println("Quitting. Bye.")

}
