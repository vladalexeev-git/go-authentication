package v1

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"sso/config"
)

func SetupHandlers(handler *gin.Engine, log *slog.Logger, cfg *config.Config) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	handler.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
