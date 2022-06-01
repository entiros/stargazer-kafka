package stargazer_kafka

import (
	"github.com/spf13/viper"
)

// Config is the configuration file struct
type Config struct {
	Starlify struct {
		BaseUrl  string `yaml:"baseUrl"`
		ApiKey   string `yaml:"apiKey"`
		AgentId  string `yaml:"agentId"`
		SystemId string `yaml:"systemId"`
	} `yaml:"starlify"`
	Kafka struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"kafka"`
}

// LoadConfig will load YAML configuration file
func LoadConfig(file string) (*Config, error) {
	viper.SetConfigFile(file)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config Config
	err = viper.Unmarshal(&config)
	return &config, nil
}
