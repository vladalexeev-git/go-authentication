package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-authentication/config"
	"go-authentication/internal/domain"
	"go-authentication/internal/service"
	"go-authentication/pkg/apperrors"
	"go-authentication/pkg/utils"
	"log/slog"
	"net/http"
)

type accountHandler struct {
	log *slog.Logger
	cfg *config.Config

	accountService service.Account
	authService    service.Auth
}

func newAccountHandler(handler *gin.RouterGroup, log *slog.Logger, cfg *config.Config, accService service.Account, sessionService service.Session, authService service.Auth) {
	h := &accountHandler{log: log, cfg: cfg, accountService: accService, authService: authService}

	g := handler.Group("/account")

	authenticated := g.Group("/", sessionMiddleware(log, cfg, sessionService))
	{
		secure := authenticated.Group("/", tokenMiddleware(log, cfg, authService))
		{
			secure.DELETE("", h.delete)
		}

		authenticated.GET("", h.get)
	}

	g.POST("", h.create)
}

//TODO: Think, maybe import domain models in this layer is a bad practice?
//TODO: Think how to organize information logs properly

func (h *accountHandler) create(c *gin.Context) {
	const op = "api.create"
	l := h.log.With(slog.String(utils.Operation, op))
	var r accountCreateRequest

	err := c.BindJSON(&r)
	if err != nil {
		l.Error("can't unmarshal account request", slog.String("error", err.Error()))

		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{Error: apperrors.ErrorValidate.Error()})
		return
	}

	account := domain.Account{Email: r.Email, Username: r.Username, Password: r.Password}

	_, err = h.accountService.Create(c.Request.Context(), account)
	if err != nil {
		if errors.Is(err, apperrors.ErrorAccountAlreadyExists) {
			h.log.Warn("account already exists",
				slog.String(utils.Operation, op),
				slog.String("error", err.Error()))

			c.AbortWithStatusJSON(http.StatusConflict, errorResponse{Error: apperrors.ErrorAccountAlreadyExists.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *accountHandler) get(c *gin.Context) {
	const op = "api.get"
	l := h.log.With(slog.String(utils.Operation, op))

	aid, err := getAccountID(c)
	if err != nil {
		l.Error("can't get account id", slog.String("error", err.Error()))

		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	acc, err := h.accountService.GetByID(c.Request.Context(), aid)
	if err != nil {
		if errors.Is(err, apperrors.ErrorAccountNotFound) {
			l.Warn("account not found", slog.String("account id", aid), slog.String("error", err.Error()))

			c.AbortWithStatusJSON(http.StatusNotFound, errorResponse{Error: apperrors.ErrorAccountNotFound.Error()})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, acc)
}

func (h *accountHandler) delete(c *gin.Context) {
	const op = "api.delete"
	l := h.log.With(slog.String(utils.Operation, op))

	aid, err := getAccountID(c)
	if err != nil {
		l.Error("can't get account id", slog.String("error", err.Error()))

		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	l.Debug("deleting acc by id", slog.String("aid", aid))

	err = h.accountService.Delete(c.Request.Context(), aid)
	if err != nil {
		l.Error("can't delete account", slog.String("error", err.Error()))

		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "account was deleted",
	})
}
