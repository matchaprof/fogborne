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

func handlePlayerConnection(playerID string) {
	// Start a new player session
	session := logging.StartPlayerSession(playerID)
	logging.WithFields(logrus.Fields{
		"context": session,
		"status":  "connecting",
	}).Info("Player connection initiated")

	// Player performs an action
	action := session.StartGameAction("move")
	logging.WithFields(logrus.Fields{
		"context": action,
		"x":       10,
		"y":       20,
	}).Info("Player movement started")

	// Action completes
	logging.WithFields(logrus.Fields{
		"context": action,
		"status":  "completed",
	}).Info("Player movement completed")
}

func main() {
	cfg, err := config.LoadConfig(DEV)
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Initialize logger
	if err := logging.InitLogger(&cfg.Logging); err != nil {
		log.Fatalf("Logger initialization error: %v", err)
	}

	logging.LogTitle("Welcome to Fogborne: Terminal Horizons!", logging.Logger.Info)

	logging.LogSection("Fogborne Server Initialization", logging.Logger.Info)

	// Testing logger
	logging.LogSubSection("Configuration Loaded", logging.Logger.Info)
	logging.WithFields(logrus.Fields{
		"map_width":  cfg.Game.MapSize.Width,
		"map_height": cfg.Game.MapSize.Height,
		"tick_rate":  cfg.Game.TickRate,
		"port":       cfg.Game.ServerConfig.Port,
	}).Info("[ -- Game Configuration Initialized -- ]")
	// logging.WithFields(logrus.Fields{
	// 	"map_width":  cfg.Game.MapSize.Width,
	// 	"map_height": cfg.Game.MapSize.Height,
	// 	"tick_rate":  cfg.Game.TickRate,
	// 	"port":       cfg.Game.ServerConfig.Port,
	// }).Debug("[ -- Game Configuration Initialized -- ]")
	// logging.WithFields(logrus.Fields{
	// 	"map_width":  cfg.Game.MapSize.Width,
	// 	"map_height": cfg.Game.MapSize.Height,
	// 	"tick_rate":  cfg.Game.TickRate,
	// 	"port":       cfg.Game.ServerConfig.Port,
	// }).Warn("[ -- Game Configuration Initialized -- ]")
	// logging.WithFields(logrus.Fields{
	// 	"map_width":  cfg.Game.MapSize.Width,
	// 	"map_height": cfg.Game.MapSize.Height,
	// 	"tick_rate":  cfg.Game.TickRate,
	// 	"port":       cfg.Game.ServerConfig.Port,
	// }).Error("[ -- Game Configuration Initialized -- ]")
	// logging.WithFields(logrus.Fields{
	// 	"map_width":  cfg.Game.MapSize.Width,
	// 	"map_height": cfg.Game.MapSize.Height,
	// 	"tick_rate":  cfg.Game.TickRate,
	// 	"port":       cfg.Game.ServerConfig.Port,
	// }).Fatal("[ -- Game Configuration Initialized -- ]")

	handlePlayerConnection("5003122201")
}
