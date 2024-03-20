package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"sso/config"
	"sso/internal/service"
	"sso/pkg/apperrors"
	"sso/pkg/utils"
)

// todo: pass config here or not?
func sessionMiddleware(log *slog.Logger, cfg *config.Config, s service.Session) gin.HandlerFunc {
	const op = "sessionMiddleware"
	l := log.With(slog.String(utils.Operation, op))

	return func(c *gin.Context) {
		sid, err := c.Cookie(cfg.CookieKey)
		if err != nil {
			l.Error("session id is not passed", slog.String("error", err.Error())) //todo should we log?

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
