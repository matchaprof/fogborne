package main

import (
	"log"

	"github.com/matchaprof/fogborne/internal/core/config"
	"github.com/matchaprof/fogborne/internal/core/logging"
	render "github.com/matchaprof/fogborne/internal/render/ascii"
	"github.com/matchaprof/fogborne/internal/terminal"
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

	var (
		MAP_WIDTH   = cfg.Game.MapSize.Width
		MAP_HEIGHT  = cfg.Game.MapSize.Height
		TICK_RATE   = cfg.Game.TickRate
		SERVER_PORT = cfg.Game.ServerConfig.Port
	)

	// Testing logger
	logging.LogSubSection("Configuration Loaded", logging.Logger.Info)
	logging.WithFields(logrus.Fields{
		"map_width":  MAP_WIDTH,
		"map_height": MAP_HEIGHT,
		"tick_rate":  TICK_RATE,
		"port":       SERVER_PORT,
	}).Info(" .•( Game Configuration Initialized )•.")

	// Creating test map
	logging.WithFields(logrus.Fields{
		"map_width":  MAP_WIDTH,
		"map_height": MAP_HEIGHT,
	}).Infof(" .•( Map Settings being used by NewGameMap %dx%d )•.", MAP_WIDTH, MAP_HEIGHT)
	gameMap := render.NewGameMap(MAP_WIDTH, MAP_HEIGHT)

	// Checking player's terminal size
	termWidth, termHeight, err := terminal.GetTerminalSize()
	if err != nil {
		logging.Logger.Error("Error detecting terminal size\n")
		return
	}

	// Calculate offset in order to center the map in the player's terminal
	offsetX := (termWidth - MAP_WIDTH) / 2
	offsetY := (termHeight - MAP_HEIGHT) / 2
	if offsetX < 0 || offsetY < 0 {
		MAP_WIDTH = termWidth - 1
		MAP_HEIGHT = termHeight - 1
	}

	// Add walls
	for x := 0; x < MAP_WIDTH; x++ {
		gameMap.Tiles[0][x] = render.CeilingTile
		gameMap.Tiles[gameMap.Height-1][x] = render.LowerTile
	}

	for y := 0; y < MAP_HEIGHT; y++ {
		gameMap.Tiles[y][0] = render.WallTile
		gameMap.Tiles[y][gameMap.Width-1] = render.WallTile
	}

	// Draw the map
	gameMap.Draw(offsetX, offsetY)
}
