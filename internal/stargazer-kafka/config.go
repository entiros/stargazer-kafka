package stargazer_kafka

import (
	"github.com/spf13/viper"
	"os"
	"strings"
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
		Type string `yaml:"type"`
		Auth struct {
			OAuth struct {
				Token string `yaml:"token"`
			} `yaml:"oauth"`
		} `yaml:"auth"`
	} `yaml:"kafka"`
}

// LoadConfig will load properties from YAML configuration file or environment variables
func LoadConfig() (*Config, error) {
	// Default Starlify properties
	viper.SetDefault("starlify.baseUrl", "https://api.starlify.com/hypermedia")
	viper.SetDefault("starlify.apiKey", "")
	viper.SetDefault("starlify.systemId", "")
	viper.SetDefault("starlify.agentId", "")

	// Default Kafka properties
	viper.SetDefault("kafka.host", "127.0.0.1:9092")
	viper.SetDefault("kafka.type", "msk")
	viper.SetDefault("kafka.oauth.token", "")

	// Override properties with upper case environment variable of property name with . replaced with _
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Load configuration file
	if len(os.Args) > 1 {
		viper.SetConfigFile(os.Args[1])
		viper.SetConfigType("yaml")

		err := viper.ReadInConfig()
		if err != nil {
			return nil, err
		}
	}

	//for _, k := range viper.AllKeys() {
	//	v := viper.GetString(k)
	//	log.Printf("LOG : %s =  %s", k, v)
	//}

	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
