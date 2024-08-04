package main

import (
	"context"
	"fmt"
	"github.com/civilrights3/go-derek-go/internal/config"
	"github.com/civilrights3/go-derek-go/internal/multiworld"
	"gopkg.in/yaml.v3"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cfg, err := readConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	sigChan := make(chan os.Signal, 1)
	// catch SIGETRM or SIGINTERRUPT
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	// init the adapter for archipelago
	arch := multiworld.NewArchipelagoClient(cfg.Multiworld)
	err = arch.Start(ctx, cfg.Multiworld.World.Server, cfg.Multiworld.World.Port, cfg.Multiworld.World.Slot)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		arch.Disconnect(ctx)
		fmt.Println("Disconnected...")
		cancel()
	}()

	_, err = arch.Read(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	// init adapter for discord

	// build core and pass adapters

	fmt.Println("Started...")
	<-sigChan
	fmt.Println("Closing...")
	//cancel()
}

const (
	confLoc    = "config" //TODO make parameter with default location
	configName = "config.yaml"
	keyName    = "key"
)

func readConfig() (config.Config, error) {
	cfg := config.NewDefaultConfig()

	b, err := os.ReadFile(filepath.Join(confLoc, configName))
	if err != nil {
		return cfg, fmt.Errorf("unable to read config file: %w", err)
	}

	err = yaml.Unmarshal(b, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("unable to unmarshal config file: %w", err)
	}

	k, err := os.ReadFile(filepath.Join(confLoc, keyName))
	if err != nil {
		return cfg, fmt.Errorf("unable to api key file: %w", err)
	}

	cfg.Chat.Key = string(k)

	return cfg, nil
}
