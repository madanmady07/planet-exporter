package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"planet-exporter/cmd/internal"
	"planet-exporter/collector"

	log "github.com/sirupsen/logrus"
)

var (
	version            string
	showVersionAndExit bool
)

func main() {
	var config internal.Config

	// Main
	flag.StringVar(&config.ListenAddress, "listen-address", "0.0.0.0:19100", "Address to which exporter will bind its HTTP interface")
	flag.StringVar(&config.LogLevel, "log-level", "info", "Log level")
	flag.BoolVar(&showVersionAndExit, "version", false, "Show version and exit")

	// Collector tasks
	flag.StringVar(&config.TaskInterval, "task-interval", "7s", "Interval between collection of expensive data into memory")
	flag.StringVar(&config.DarkstatAddr, "darkstat-addr", "http://localhost:51666/metrics", "Darkstat target address")
	flag.StringVar(&config.InventoryAddr, "inventory-addr", "http://172.21.44.249:18081/inventory.json", "Inventory target address")

	flag.Parse()

	if showVersionAndExit {
		fmt.Printf("planet-exporter %v\n", version)
		os.Exit(0)
	}

	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: false,
	})
	logLevel, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		log.Fatalf("Failed to parse log level: %v", err)
	}
	log.SetLevel(logLevel)

	log.Infof("Planet Exporter %v", version)
	log.Infof("Initialize log with level %v", config.LogLevel)

	ctx := context.Background()

	log.Info("Initialize main service")
	collector, err := collector.NewPlanetCollector()
	if err != nil {
		log.Fatalf("Failed to initialize planet collector: %v", err)
	}

	log.Info("Initialize main service")
	svc := internal.New(config, collector)
	if err := svc.Run(ctx); err != nil {
		log.Fatalf("Main service exit with error: %v", err)
	}

	log.Info("Main service exit successfully")
}
