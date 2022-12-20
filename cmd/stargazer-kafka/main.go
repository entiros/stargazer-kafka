package main

import (
	"context"
	"fmt"
	"github.com/entiros/stargazer-kafka/internal/config"
	"github.com/entiros/stargazer-kafka/internal/system"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Start with config file name of config directory")
		os.Exit(1)
	}

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	info, err := os.Stat(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	if info.IsDir() {
		configFiles, err := config.GetConfigs(info.Name())
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			err := config.WatchConfigs(info.Name(), func(newFiles []string) {
				log.Printf("Config files: %v", newFiles)

				addMe, deleteMe := config.Diff(newFiles, configFiles)

				for _, c := range deleteMe {
					system.DeleteSystem(ctx, c)
					config.DeleteString(c, configFiles)
				}

				for _, c := range addMe {
					system.AddSystem(ctx, c)
					configFiles = append(configFiles, c)
				}

			})
			if err != nil {
				log.Println(err)
			}
		}()

	} else {
		system.AddSystem(ctx, info.Name())
	}

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		}
	}

}
