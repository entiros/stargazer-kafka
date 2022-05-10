package stargazer_kafka

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
)

func (cfg *Config) InitCron() {

	s := gocron.NewScheduler(time.UTC)

	_, err := s.Every(5).Seconds().Do(cfg.ReadTopics)
	if err != nil {
		fmt.Printf("failed to run cron job %s", err)
	}

	_, err = s.Every(5).Seconds().Do(cfg.ReadServices)
	if err != nil {
		fmt.Printf("failed to run cron job %s", err)
	}

	s.StartBlocking()
}
