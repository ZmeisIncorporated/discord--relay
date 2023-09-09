package main

import (
	"context"
	"flag"

	"os"
	"os/signal"
	"syscall"

	"log"
	
	"github.com/ZmeisIncorporated/discord--relay/internal/config"
	"github.com/ZmeisIncorporated/discord--relay/internal/forwarder"
	"github.com/ZmeisIncorporated/discord--relay/internal/listener"
)


func main() {
	var configPath = flag.String("c", "config.yaml", "specicies the path to the config")
	// TODO
	_, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	flag.Parse()

	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		log.Fatalf("config doesn't exist: %s", err)
	}
	log.Printf("use config: %s", *configPath)

	cfg, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatalf("Error getting config: %s", err)
	}

	// Connect to the destination for all messages
	f, err := forwarder.NewForwarder(cfg.Webhooks)
	if err != nil {
		log.Fatalf("Error while creating Forwarder: %s", err)
	}

	// Test that we can send to this correctly
	err = f.Send("[log]", "Forwarder Connected")
	if err != nil {
		log.Fatalf("Error while sending to forwarders: %s", err)
	}

	// Open up pidgin logs listener
	pidgin, err := listener.NewPidginListener(f, cfg.Logs)
	if err != nil {
		log.Fatalf("Error while sending to log: %s", err)
	}
	pidgin.Run()
	log.Printf("Relay is now running.  Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	oscall := <-sc
	log.Printf("Shutting down from: %+v", oscall)

}
