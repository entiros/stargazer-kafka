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
		log.Fatal("Start with config file name or directory with config files")
	}

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	info, err := os.Stat(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

loop:
	for {
		next, hasNext := getSystems(info, ctx)

		for hasNext() {
			system, err := next()
			if err != nil {
				log.Printf("Failed to sync. %v", err)
				continue
			}

			err = system.PingStarlify(ctx)
			if err != nil {
				log.Printf("failed to ping starlify. %v", err)
			}
			err = system.SyncTopics(ctx)
			if err != nil {
				log.Printf("failed to sync topics for %s, %v ", system.Name(), err)
			}

			time.Sleep(3 * time.Second)

		}

		select {
		case <-ctx.Done():
			break loop
		case <-time.After(10 * time.Second):
		}
	}

}

func getSystems(info os.FileInfo, ctx context.Context) (next func() (*system.System, error), hasNext func() bool) {

	if info.IsDir() {
		return GetSystems(ctx, info.Name())

	} else {
		return GetSystem(ctx, info.Name())
	}
}

func GetSystem(ctx context.Context, dir string) (next func() (*system.System, error), hasNext func() bool) {

	var i int
	return func() (*system.System, error) {
			s, err := system.NewSystem(ctx, dir)
			i++
			return s, err
		}, func() bool {
			return i < 1
		}

}

func GetSystems(ctx context.Context, dir string) (next func() (*system.System, error), hasNext func() bool) {

	configFiles, err := config.GetConfigs(dir)
	if err != nil {
		return nil, func() bool {
			return false
		}
	}

	var i int

	next = func() (*system.System, error) {
		sys, err := system.NewSystem(ctx, configFiles[i])
		i++
		return sys, err
	}

	hasNext = func() bool {
		return i < len(configFiles)
	}

	return
}

func GetSystemsZ(ctx context.Context, dir string) ([]*system.System, error) {

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
