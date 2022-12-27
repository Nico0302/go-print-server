package main

import (
	"log"

	"github.com/gresio/cloudprint/config"
	"github.com/gresio/cloudprint/internal/app"
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
