package main

import (
	"log"

	"github.com/matchaprof/fogborne/internal/core/config"
	"github.com/matchaprof/fogborne/internal/core/logging"
	"github.com/sirupsen/logrus"
)

const (
	DEV  = "dev"
	PROD = "prod"
	TEST = "test"
)

func main() {
	cfg, err := config.LoadConfig(DEV)
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Initialize logger
	if err := logging.InitLogger(&cfg.Logging); err != nil {
		log.Fatalf("Logger initialization error: %v", err)
	}

	logging.LogTitle("Welcome to Fogborne: An ASCII-based Survival Game", logging.Logger.Info)

	logging.LogSection("Fogborne Server Initialization", logging.Logger.Info)

	// Testing logger
	logging.LogSubSection("Configuration Loaded", logging.Logger.Info)
	logging.WithFields(logrus.Fields{
		"map_width":  cfg.Game.MapSize.Width,
		"map_height": cfg.Game.MapSize.Height,
		"tick_rate":  cfg.Game.TickRate,
		"port":       cfg.Game.ServerConfig.Port,
	}).Info(" .•( Game Configuration Initialized )•.")
}
