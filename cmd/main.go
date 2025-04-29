package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/civilrights3/go-derek-go/internal/chat"
	"github.com/civilrights3/go-derek-go/internal/config"
	"github.com/civilrights3/go-derek-go/internal/multiworld"
	"github.com/civilrights3/go-derek-go/internal/queue"
	"github.com/civilrights3/go-derek-go/test/mock"
	"gopkg.in/yaml.v3"
)

var mockArchi = flag.Bool("mockarchi", false, "use mock archipelago messages")

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cfg, err := readConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Loaded config")

	sigChan := make(chan os.Signal, 1)
	// catch SIGETRM or SIGINTERRUPT
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	// init adapter for discord
	discordClient, err := chat.NewDiscordClient(cfg.Chat)
	if err != nil {
		panic(fmt.Sprintf("error creating discord connection: %s\n", err))
	}

	err = discordClient.Connect()
	if err != nil {
		panic(fmt.Sprintf("cannot start discord connection: %s\n", err))
	}

	// start message queue
	queue.StartMessageQueue()
	//queue.Queue.RegisterMessageListener(queue.Queue.TestHandler)
	queue.Queue.RegisterMessageListener(discordClient.SendMessage)

	// init the adapter for archipelago
	if *mockArchi {
		fmt.Println("Sending mocked messages")
		mock.SendTestMessages(ctx)
		fmt.Println("Sent mocked messages")
	} else {
		arch, err := multiworld.NewArchipelagoClient(cfg.Multiworld)
		if err != nil {
			panic(fmt.Sprintf("cannot start multiworld connection: %s\n", err))
		}
		fmt.Println("Starting multiworld connection")
		arch.Start(ctx, cfg.Multiworld.World.Server, cfg.Multiworld.World.Port, cfg.Multiworld.World.Slot)
		fmt.Println("Multiworld connected")
	}

	// build core and pass adapters

	fmt.Println("Started...")
	select {
	case <-sigChan:
		break
	}

	fmt.Println("Closing...")
	err = discordClient.Disconnect()
	if err != nil {
		fmt.Printf("could not disconnect from discord: %s\n", err)
	}

	cancel()
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
