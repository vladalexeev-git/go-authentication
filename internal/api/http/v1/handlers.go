package v1

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"sso/config"
	"sso/internal/service"
)

const apiPath = "/v1"

func SetupHandlers(handler *gin.Engine, log *slog.Logger, cfg *config.Config, accService service.Account, authService service.Auth) {

	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	//handler.Static(fmt.Sprintf("%s/swagger/", apiPath), "third_party/swaggerui")

	handler.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	h := handler.Group(apiPath)

	{
		newAccountHandler(h, log, cfg, accService)
		newAuthHandler(h, log, cfg, authService)
	}

}
