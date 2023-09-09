package main

import (
	"log"
	"flag"
	"fmt"

	"os"
	"os/signal"
	"syscall"

	"github.com/ZmeisIncorporated/discord--relay/internal/config"
	"github.com/ZmeisIncorporated/discord--relay/internal/forwarder"
	"github.com/ZmeisIncorporated/discord--relay/internal/listener"
	"github.com/ZmeisIncorporated/discord--relay/internal/monitor"
)


func main() {
	var configPath = flag.String("c", "config.yaml", "specicies the path to the config")

	// Parse up the command line flags
	flag.Parse()

	// Make sure the config exists
	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		log.Printf("config doesnt exist: %s", err)
		return
	}

	//Load Config
	cfg, err := config.NewConfig(*configPath)
	if err != nil {
		log.Printf("Error getting config: %s", err)
		return
	}

	// Connect to the destination for all messages
	f, err := forwarder.NewForwarder(cfg.Webhooks, cfg.Admhooks)
	if err != nil {
		log.Printf("Error while creating Forwarder: %s", err)
		return
	}
	f.Start()
	log.Println("Forwarder Connected")
	f.AdmSend("[main]", "Forwarder Connected")

	// Start finch monitoring
	m := monitor.NewMonitor(f)
	m.Start()
	log.Println("Monitor Started")
	f.AdmSend("[main]", "Monitor Started")

	// Open up pidgin logs listener
	pidgin := listener.NewPidginListener(f, cfg.Logs)
	pidgin.Start()
	log.Println("PidginListener Started")
	f.AdmSend("[main]", "PidginListener Started")

	log.Println("Relay is now running.  Press CTRL-C to exit.")
	f.AdmSend("[main]", "Relay is now running.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	oscall := <-sc

	msg := fmt.Sprintf("Shutting down from: %+v", oscall)
	log.Printf(msg)
}
