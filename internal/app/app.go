package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/nico0302/go-print-server/config"
	"github.com/nico0302/go-print-server/internal/auth"
	"github.com/nico0302/go-print-server/internal/controller"
	"github.com/nico0302/go-print-server/internal/fetcher"
	"github.com/nico0302/go-print-server/pkg/httpserver"
	"github.com/nico0302/go-print-server/pkg/logger"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
	printerCtx := controller.NewPrinterContext(*cfg)
	fetcher := fetcher.New()
	auth := auth.New(cfg.HTTP.Users)

	// HTTP Server
	if cfg.Log.Level != string(logger.DebugLevel) {
		gin.SetMode(gin.ReleaseMode)
	}
	handler := gin.New()
	controller.NewRouter(handler, l, printerCtx, fetcher, auth)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	l.Info(fmt.Sprintf("Running print server on port %s", cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err := httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
