package config

import (
	"fmt"
	"github.com/entiros/stargazer-kafka/internal/log"
	"github.com/spf13/viper"
	"os"
	"strings"
)

// Config is the configuration file struct
type Config struct {
	Sync struct {
		Direction string `json:"direction"`
	} `yaml:"sync"`

	Starlify struct {
		BaseUrl      string `yaml:"baseUrl"`
		ApiKey       string `yaml:"apiKey"`
		AgentId      string `yaml:"agentId"`
		MiddlewareId string `yaml:"middlewareId"`
	} `yaml:"starlify"`

	Kafka struct {
		BootstrapServers []string `yaml:"bootstrapServers"`
		Auth             struct {
			OAuth struct {
				Token string `yaml:"token"`
			} `yaml:"oauth"`
			IAM struct {
				Key    string `yaml:"key"`
				Secret string `yaml:"secret"`
			} `yaml:"iam"`
			Plain struct {
				Username string `yaml:"username"`
				Password string `yaml:"password"`
			} `yaml:"plain"`
		} `yaml:"auth"`
	} `yaml:"kafka"`
}

// LoadConfig will load properties from YAML configuration file or environment variables
func LoadConfig(configFile string) (*Config, error) {

	viper.SetDefault("sync.direction", "starlify_to_kafka")

	// Default Starlify properties
	viper.SetDefault("starlify.baseUrl", "https://api.starlify.com/hypermedia")
	viper.SetDefault("starlify.apiKey", "")
	viper.SetDefault("starlify.middlewareId", "")
	viper.SetDefault("starlify.agentId", "")

	// Default Kafka properties
	viper.SetDefault("kafka.bootstrapServers", []string{"127.0.0.1:9092"})
	viper.SetDefault("kafka.auth.oauth.token", "")
	viper.SetDefault("kafka.auth.plain.username", "")
	viper.SetDefault("kafka.auth.plain.password", "")
	viper.SetDefault("kafka.auth.iam.secret", "")
	viper.SetDefault("kafka.auth.iam.key", "")

	// Override properties with upper case environment variable of property name with . replaced with _
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Load configuration file
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to load config for %s. %v", configFile, err)
	}

	return &config, nil
}

func GetConfigs(dir string) ([]string, error) {

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && entry.Type().IsRegular() && (strings.HasSuffix(entry.Name(), "yml") || strings.HasSuffix(entry.Name(), "yaml")) {
			files = append(files, dir+string(os.PathSeparator)+entry.Name())
		}
	}

	log.Logger.Debugf("Found %d config files in %s", len(files), dir)
	for _, file := range files {
		log.Logger.Debugf("config file: %s", file)
	}

	return files, nil
}

func DeleteString(deleteMe string, list []string) {

	for i, d := range list {
		if d == deleteMe {
			list = append(list[:i], list[i+1:]...)
		}
	}

}

func contains(s string, in []string) bool {

	for _, v := range in {
		if v == s {
			return true
		}
	}
	return false

}

func Diff(newFiles, oldFiles []string) ([]string, []string) {

	var deleteMe []string
	var addMe []string
	for _, f := range oldFiles {
		if !contains(f, newFiles) {
			deleteMe = append(deleteMe, f)
		}
	}

	for _, f := range newFiles {
		if !contains(f, oldFiles) {
			addMe = append(addMe, f)
		}
	}
	return addMe, deleteMe

}
