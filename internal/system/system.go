package system

import (
	"context"
	"fmt"
	"github.com/entiros/stargazer-kafka/internal/config"
	"github.com/entiros/stargazer-kafka/internal/kafka"
	stargazerkafka "github.com/entiros/stargazer-kafka/internal/stargazer-kafka"
	"github.com/entiros/stargazer-kafka/internal/starlify"
)

type System struct {
	cfg    *config.Config
	cancel context.CancelFunc
	ctx    context.Context
	file   string
	ks     *stargazerkafka.KafkaTopicsToStarlify
}

func (s *System) Name() string {
	return s.file
}

var systems map[string]*System

func init() {
	systems = make(map[string]*System)
}

func NewSystem(ctx context.Context, c string) (*System, error) {

	cfg, err := config.LoadConfig(c)
	if err != nil {
		return nil, err
	}

	s := &System{
		cfg:  cfg,
		file: c,
	}
	err = s.init(ctx)
	if err != nil {
		return nil, err
	}
	return s, nil

}

func (s *System) init(ctx context.Context) error {

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
		return fmt.Errorf("failed to initialize system %s. %v", s.file, err)
	}

	s.ks = kafkaTopicsToStarlify

	return nil
}

func (s *System) SyncTopics(ctx context.Context) error {
	return s.ks.SyncTopics(ctx)
}

func (s *System) PingStarlify(ctx context.Context) error {
	return s.ks.Ping(ctx)
}
