package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/nico0302/go-print-server/internal/auth"
	"github.com/nico0302/go-print-server/internal/fetcher"
	"github.com/nico0302/go-print-server/pkg/logger"
)

func NewRouter(handler *gin.Engine, l logger.Interface, c *PrinterContext, f *fetcher.Fetcher, a *auth.Auth) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())
	handler.Use(CORSMiddleware())

	// Routers
	h := handler.Group("/v1", a.GetHandlerFunc())
	{
		newPrintRoutes(h, c, f, l)
	}
}
