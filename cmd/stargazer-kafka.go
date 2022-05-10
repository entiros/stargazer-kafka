package main

import (
	"fmt"
	stargazerkafka "github.com/entiros/stargazer-kafka"
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
	defer f.Close()

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

	// Read config file and populate Config
	readConfigFile(cfg)
	cfg.InitCron()
}
