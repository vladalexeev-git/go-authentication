package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"sso/config"
	"sso/internal/service"
	"sso/pkg/apperrors"
)

type authHandler struct {
	log *slog.Logger
	cfg *config.Config

	auth service.Auth
}

func newAuthHandler(handler *gin.RouterGroup, log *slog.Logger, cfg *config.Config, auth service.Auth) {
	h := &authHandler{log: log, cfg: cfg, auth: auth}

	g := handler.Group("/login")

	{
		g.POST("/", h.login)

	}
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *authHandler) login(c *gin.Context) {
	const op = "api.login"
	l := h.log.With(slog.String("operation", op))
	var r loginRequest

	err := c.Bind(&r)
	if err != nil {
		l.Error("can't unmarshal login request", slog.String("error", err.Error()))

		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: apperrors.ErrorValidate.Error()})
		return
	}

	l.Debug("login request",
		slog.Any("request", r),
		slog.String("email", r.Email),
		slog.String("user-agent", c.Request.UserAgent()),
		slog.String("ip", c.ClientIP()))

	s, err := h.auth.EmailLogin(
		c.Request.Context(),
		r.Email,
		r.Password,
		service.Device{
			UserAgent: c.Request.UserAgent(),
			IP:        c.ClientIP(),
		})
	if err != nil {
		if errors.Is(err, apperrors.ErrorAccountNotFound) ||
			errors.Is(err, apperrors.ErrorAccountWrongPassword) {
			l.Warn("email or password incorrect")
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: apperrors.ErrorLoginOrPasswordIncorrect.Error()})
		}
		l.Warn("cannot login", slog.String("error", err.Error()))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.SetCookie(
		h.cfg.Session.CookieKey,
		s.ID,
		s.TTL,
		h.cfg.Session.CookiePath,
		h.cfg.Session.CookieDomain,
		h.cfg.Session.CookieSecure,
		h.cfg.Session.CookieHttpOnly,
	)
	c.AbortWithStatus(http.StatusOK)

}

//todo logout
