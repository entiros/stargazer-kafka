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
	"time"
)

func main() {

	fmt.Println("Starting Stargazer")
	if len(os.Args) < 2 {
		log.Fatal("Start with config file name of config directory")
	}

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	info, err := os.Stat(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

loop:
	for {
		systems := getSystems(info, ctx)

		for _, system := range systems {
			system.PingStarlify(ctx)
			system.SyncTopics(ctx)
		}
		select {
		case <-ctx.Done():
			break loop
		case <-time.After(10 * time.Second):
		}
	}

}

func getSystems(info os.FileInfo, ctx context.Context) []*system.System {

	var systems []*system.System
	var err error

	if info.IsDir() {
		systems, err = GetSystems(ctx, info.Name())
		if err != nil {
			log.Fatal(err)
		}
	} else {
		sys, e := system.NewSystem(ctx, info.Name())
		if e != nil {
			log.Fatal(e)
		}
		systems = append(systems, sys)
	}
	return systems
}

func GetSystems(ctx context.Context, dir string) ([]*system.System, error) {

	configFiles, err := config.GetConfigs(dir)
	if err != nil {
		return nil, err
	}

	var systems []*system.System

	for _, config := range configFiles {
		sys, err := system.NewSystem(ctx, config)
		if err != nil {
			systems = append(systems, sys)
		}
	}

	return systems, nil

}

/*
func WatchDir2(ctx context.Context, dir string) error {

	configFiles, err := config.GetConfigs(dir)
	if err != nil {
		return err
	}

	go func() {
		err := config.WatchConfigs(dir, func(newFiles []string) {
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

	return nil
}

func WatchDir2(ctx context.Context, dir string) error {

	configFiles, err := config.GetConfigs(dir)
	if err != nil {
		return err
	}

	go func() {
		err := config.WatchConfigs(dir, func(newFiles []string) {
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

	return nil
}
*/
