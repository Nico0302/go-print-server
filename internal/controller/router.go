package controller

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gresio/cloudprint/internal/auth"
	"github.com/gresio/cloudprint/internal/fetcher"
	"github.com/gresio/cloudprint/pkg/logger"
)

func NewRouter(handler *gin.Engine, l logger.Interface, c *PrinterContext, f *fetcher.Fetcher, a *auth.Auth) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())
	handler.Use(cors.Default())

	// Routers
	h := handler.Group("/v1", a.GetHandlerFunc())
	{
		newPrintRoutes(h, c, f, l)
	}
}
