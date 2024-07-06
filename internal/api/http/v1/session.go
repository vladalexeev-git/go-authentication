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

type sessionHandler struct {
	l    *slog.Logger
	cfg  *config.Config
	sess service.Session
}

func newSessionHandler(
	handler *gin.RouterGroup,
	l *slog.Logger,
	cfg *config.Config,
	sess service.Session,
	auth service.Auth) {

	h := &sessionHandler{l: l, cfg: cfg, sess: sess}
	g := handler.Group("/session")
	{
		authenticated := g.Group("/", sessionMiddleware(l, cfg, sess))
		{
			secure := authenticated.Group("/", tokenMiddleware(l, cfg, auth))
			{
				secure.DELETE(":sessionID", h.terminate)
				secure.DELETE("", h.terminateAll)
			}
			authenticated.GET(":sessionID", h.get)
		}
	}
}

func (h *sessionHandler) terminate(c *gin.Context) {
	const op = "api.terminate"
	l := h.l.With(slog.String(utils.Operation, op))

	// get the current session id
	curSid, err := getSessionID(c)
	if err != nil {
		l.Error("can't get session id", slog.String("error", err.Error()))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = h.sess.Terminate(c.Request.Context(), curSid, c.Param("sessionID"))
	if err != nil {
		if errors.Is(err, apperrors.ErrorCurrentSessionTerminating) {
			c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{Error: apperrors.ErrorCurrentSessionTerminating.Error()})
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
	return
}

func (h *sessionHandler) terminateAll(c *gin.Context) {
	const op = "api.terminateAll"
	l := h.l.With(slog.String(utils.Operation, op))

	// get the current session id
	curSid, err := getSessionID(c)
	if err != nil {
		l.Error("can't get session id", slog.String("error", err.Error()))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// and acc id
	aid, err := getAccountID(c)
	if err != nil {
		l.Error("can't get account id", slog.String("error", err.Error()))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = h.sess.TerminateAll(c.Request.Context(), curSid, aid)
	if err != nil {
		l.Error("can't terminate sessions", slog.String("error", err.Error()))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	l.Debug("all sessions of account successfully terminated", slog.String("account id", aid))

	c.Status(http.StatusNoContent)
	return
}

func (h *sessionHandler) get(c *gin.Context) {
	const op = "api.get"
	l := h.l.With(slog.String(utils.Operation, op))

	aid, err := getAccountID(c)
	if err != nil {
		l.Error("can't get account id", slog.String("error", err.Error()))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	sessions, err := h.sess.GetAll(c.Request.Context(), aid)
	if err != nil {
		l.Error("can't get sessions", slog.String("error", err.Error()))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, sessions)

}
