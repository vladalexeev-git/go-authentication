package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-authentication/config"
	"go-authentication/internal/apperrors"
	"go-authentication/internal/service"
	"go-authentication/pkg/utils"
	"log/slog"
	"net/http"
)

func sessionMiddleware(log *slog.Logger, cfg *config.Config, s service.Session) gin.HandlerFunc {
	const op = "sessionMiddleware"
	l := log.With(slog.String(utils.Operation, op))

	return func(c *gin.Context) {
		sid, err := c.Cookie(cfg.CookieKey)
		if err != nil {
			l.Warn("session id is not passed", slog.String("error", err.Error()))

			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		l.Debug("got session key from cookie", slog.String("sid", sid))

		session, err := s.Get(c.Request.Context(), sid)
		if err != nil {
			l.Warn("session not found", slog.String("error", err.Error()))

			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		l.Debug("got account id from session service", slog.String("aid", session.AccountID))

		d := service.Device{
			UserAgent: c.Request.UserAgent(), //c.Request.Header.Get("User-Agent")
			IP:        c.ClientIP(),
		}

		if session.IP != d.IP || session.UserAgent != d.UserAgent {

			l.Warn("ip or user agent is different", slog.String("error", apperrors.ErrorSessionDeviceMismatch.Error()))

			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Set("sid", session.ID)
		c.Set("aid", session.AccountID)
		c.Next()
	}
}

func setCSRFTokenMiddleware(log *slog.Logger, cfg *config.Config) gin.HandlerFunc {
	const op = "setCSRFTokenMiddleware"
	l := log.With(slog.String(utils.Operation, op))

	return func(c *gin.Context) {

		c.Next()

		t, err := utils.UniqueString(32)
		if err != nil {
			l.Error("can't generate csrf token", slog.String("error", err.Error()))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		l.Info("csrf token middleware", slog.String("token", t), slog.String("header key", cfg.CSRFHeaderKey), slog.String("cookie key", cfg.CSRFCookieKey))

		c.Header(cfg.CSRFHeaderKey, t)
		c.SetCookie(
			cfg.CSRFCookieKey,
			t,
			int(cfg.CSRFttl.Seconds()),
			apiPath,
			cfg.CookieDomain,
			cfg.CookieSecure,
			cfg.CookieHttpOnly,
		)
	}
}

func csrfMiddleware(log *slog.Logger, cfg *config.Config) gin.HandlerFunc {
	const op = "csrfTokenMiddleware"
	l := log.With(slog.String(utils.Operation, op))

	return func(c *gin.Context) {
		ct, err := c.Cookie(cfg.CSRFCookieKey)
		if err != nil {
			l.Warn("csrf token is not passed", slog.String("error", err.Error()))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ht := c.GetHeader(cfg.CSRFHeaderKey)

		if ct != ht || ht == "" || ct == "" {
			l.Warn("csrf token is invalid", slog.String("cookie token", ct), slog.String("header token", ht))
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		l.Info("csrf token middleware", slog.String("passed cookie token", ct), slog.String("passed header token", ht))
		c.Next()

		t, err := utils.UniqueString(32) //todo: maybe change to uuid
		if err != nil {
			l.Error("can't generate csrf token", slog.String("error", err.Error()))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		l.Info("csrf token middleware", slog.String("new token (value)", t), slog.String("header name", cfg.CSRFHeaderKey), slog.String("cookie name", cfg.CSRFCookieKey))

		c.Header(cfg.CSRFHeaderKey, t)
		c.SetCookie(
			cfg.CSRFCookieKey,
			t,
			int(cfg.CSRFttl.Seconds()),
			apiPath,
			cfg.CookieDomain,
			cfg.CookieSecure,
			cfg.CookieHttpOnly,
		)
	}
}

func tokenMiddleware(log *slog.Logger, cfg *config.Config, a service.Auth) gin.HandlerFunc {
	const op = "tokenMiddleware"
	l := log.With(slog.String(utils.Operation, op))

	return func(c *gin.Context) {
		aid, err := getAccountID(c)
		if err != nil {
			l.Warn("account id is empty", slog.String("error", err.Error()))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		t, found := c.GetQuery("token")
		if !found || t == "" {
			l.Warn("access token is not passed")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		sub, err := a.ParseAccessToken(c.Request.Context(), t)
		if err != nil {
			l.Warn("access token is invalid", slog.String("error", err.Error()))
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		if aid != sub {
			l.Warn("access token is invalid")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}

func getAccountID(c *gin.Context) (string, error) {
	aid := c.GetString("aid")
	_, err := uuid.Parse(aid)
	if err != nil {
		return "", apperrors.ErrorContextAccountIdNotFount
	}
	return aid, nil
}

func getSessionID(c *gin.Context) (string, error) {
	sid := c.GetString("sid")
	if sid == "" {
		return "", apperrors.ErrorContextSessionNotFound
	}
	return sid, nil
}
