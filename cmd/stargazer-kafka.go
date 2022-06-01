package main

import (
	"fmt"
	stargazerkafka "github.com/entiros/stargazer-kafka"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gopkg.in/yaml.v3"
	"os"
)

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func readConfigFile(cfg *stargazerkafka.Config) {
	f, err := os.Open("config.yml")
	if err != nil {
		processError(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal().Msg("failed to close")
		}
	}(f)

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
}

func main() {
	//stargazerkafka.CreateTopic("test1")
	//stargazerkafka.CreateTopic("test2")
	//stargazerkafka.CreateTopic("test3")
	// Instantitating a empty Config object
	cfg := &stargazerkafka.Config{}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Info().Msg("parsing config.yml file")

	// Read config file and populate Config
	readConfigFile(cfg)

	log.Info().Msg("initiating cron jobs")
	cfg.InitCron()
}
