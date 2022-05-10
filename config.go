package stargazer_kafka

type Config struct {
	Starlify struct {
		ApiKey      string `yaml:"apiKey"`
		SystemID    string `yaml:"systemID"`
		ApiUrl      string `yaml:"apiUrl"`
		StargazerID string `yaml:"stargazerID"`
	} `yaml:"starlify"`
	Kafka struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"kafka"`
}
