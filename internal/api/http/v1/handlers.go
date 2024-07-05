package v1

import (
	"github.com/gin-gonic/gin"
	"go-authentication/config"
	"go-authentication/internal/service"
	"log/slog"
)

const apiPath = "/v1"

func SetupHandlers(
	handler *gin.Engine,
	log *slog.Logger,
	cfg *config.Config,
	acc service.Account,
	sess service.Session,
	auth service.Auth,
) {

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
		newAccountHandler(h, log, cfg, acc, sess, auth)
		newAuthHandler(h, log, cfg, auth, sess)
		//todo session handler
	}

}
