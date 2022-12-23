package system

import (
	"context"
	"github.com/entiros/stargazer-kafka/internal/config"
	"github.com/entiros/stargazer-kafka/internal/kafka"
	stargazerkafka "github.com/entiros/stargazer-kafka/internal/stargazer-kafka"
	"github.com/entiros/stargazer-kafka/internal/starlify"
	"log"
	"os"
)

type System struct {
	cfg    *config.Config
	cancel context.CancelFunc
	ctx    context.Context
	file   string
	ks     *stargazerkafka.KafkaTopicsToStarlify
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
	s.init(ctx)
	return s, nil

}

func (s *System) init(ctx context.Context) {

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

	s.ks = kafkaTopicsToStarlify

}

func (s *System) SyncTopics(ctx context.Context) error {

	err := s.ks.SyncTopics(ctx)
	if err != nil {
		s.ks.ReportError(ctx, err)
	}
	return err
}

func (s *System) PingStarlify(ctx context.Context) error {
	return s.ks.Ping(ctx)

}
