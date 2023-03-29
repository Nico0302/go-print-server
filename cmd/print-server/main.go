package main

import (
	"log"

	"github.com/nico0302/go-print-server/config"
	"github.com/nico0302/go-print-server/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
