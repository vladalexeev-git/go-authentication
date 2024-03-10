package v1

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"sso/config"
	"sso/internal/service"
)

type authHandler struct {
	log *slog.Logger
	cfg *config.Config

	auth service.Auth
}

func newAuthHandler(handler *gin.RouterGroup, log *slog.Logger, cfg *config.Config) {
	h := &authHandler{log: log, cfg: cfg}

	g := handler.Group("/login")

	{
		g.POST("/", h.login)

	}

}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (ah *authHandler) login(c *gin.Context) {
	const op = "api.login"
	var r loginRequest

	log := ah.log.With(slog.String("operation", op))
	err := c.Bind(&r)
	if err != nil {
		log.Error("can't unmarshal login request",
			slog.String("error", err.Error()))
	}

	//TODO: request to the auth service..
}
