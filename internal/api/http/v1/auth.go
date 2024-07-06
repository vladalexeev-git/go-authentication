package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-authentication/config"
	"go-authentication/internal/apperrors"
	"go-authentication/internal/service"
	"go-authentication/pkg/utils"
	"log/slog"
	"net/http"
)

type authHandler struct {
	l   *slog.Logger
	cfg *config.Config

	auth service.Auth
	sess service.Session
}

func newAuthHandler(handler *gin.RouterGroup, log *slog.Logger, cfg *config.Config, auth service.Auth, sess service.Session) {
	h := &authHandler{l: log, cfg: cfg, auth: auth, sess: sess}

	g := handler.Group("/auth")
	{
		authenticated := g.Group("", csrfMiddleware(log, cfg), sessionMiddleware(log, cfg, sess))

		{
			authenticated.POST("/logout", h.logout)
			authenticated.GET("/token", h.token)
		}

		g.POST("/login", h.login).Use(setCSRFTokenMiddleware(log, cfg))
	}
}

func (h *authHandler) login(c *gin.Context) {
	const op = "api.login"
	l := h.l.With(slog.String("operation", op))
	var r loginRequest

	err := c.Bind(&r)
	if err != nil {
		l.Error("can't unmarshal login request", slog.String("error", err.Error()))

		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{Error: apperrors.ErrorValidate.Error()})
		return
	}

	l.Info("login request",
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
			c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{Error: apperrors.ErrorLoginOrPasswordIncorrect.Error()})
		}
		l.Warn("cannot login", slog.String("error", err.Error()))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.SetCookie(
		h.cfg.Session.CookieKey,
		s.ID,
		s.TTL,
		apiPath,
		h.cfg.Session.CookieDomain,
		h.cfg.Session.CookieSecure,
		h.cfg.Session.CookieHttpOnly,
	)
	c.AbortWithStatus(http.StatusOK)

}

func (h *authHandler) logout(c *gin.Context) {
	const op = "api.logout"
	l := h.l.With(slog.String(utils.Operation, op))

	sid, err := getSessionID(c)
	if err != nil {
		l.Warn("can't get sid from cookie", slog.String("error", err.Error()))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err = h.auth.Logout(c.Request.Context(), sid); err != nil {
		l.Warn("can't logout", slog.String("error", err.Error()))

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.SetCookie( //todo why we set cookie anyway?
		h.cfg.Session.CookieKey,
		"",
		-1,
		apiPath,
		h.cfg.Session.CookieDomain,
		h.cfg.Session.CookieSecure,
		h.cfg.Session.CookieHttpOnly,
	)

	c.AbortWithStatus(http.StatusNoContent)
}

func (h *authHandler) token(c *gin.Context) {
	const op = "api.token"
	l := h.l.With(slog.String(utils.Operation, op))

	var r tokenRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		l.Error("can't unmarshal token request", slog.String("error", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{Error: apperrors.ErrorValidate.Error()})
		return
	}

	aid, err := getAccountID(c)
	if err != nil {
		l.Warn("can't get account id", slog.String("error", err.Error()))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	t, err := h.auth.NewAccessToken(c.Request.Context(), aid, r.Password)
	if err != nil {
		l.Error("", slog.String("error", err.Error()))
		if errors.Is(err, apperrors.ErrorAccountWrongPassword) {
			c.AbortWithStatus(http.StatusForbidden)
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, tokenResponse{AccessToken: t})
	return
}
