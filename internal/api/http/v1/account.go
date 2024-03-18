package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"sso/config"
	"sso/internal/domain"
	"sso/internal/service"
	"sso/pkg/apperrors"
	"sso/pkg/utils"
)

type accountHandler struct {
	log *slog.Logger
	cfg *config.Config

	accountService service.Account
}

func newAccountHandler(handler *gin.RouterGroup, log *slog.Logger, cfg *config.Config, accService service.Account, sessionService service.Session) {
	h := &accountHandler{log: log, cfg: cfg, accountService: accService}

	g := handler.Group("/account")
	g.POST("", h.create)

	authenticated := g.Group("/", sessionMiddleware(log, cfg, sessionService))
	{
		authenticated.GET("", h.get)
		authenticated.DELETE("", h.delete)
	}
}

//TODO: Create special errors understandable for users
//TODO: Think, maybe import domain models in this layer is a bad practice?

//TODO: add context instead of context.TODO

func (h *accountHandler) create(c *gin.Context) {
	const op = "api.create"
	l := h.log.With(slog.String(utils.Operation, op))
	var r accountCreateRequest

	err := c.BindJSON(&r)
	if err != nil {
		l.Error("can't unmarshal account request",
			slog.String(utils.Operation, op),
			slog.String("error", err.Error()))

		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: apperrors.ErrorValidate.Error()})
		return
	}

	account := domain.Account{Email: r.Email, Username: r.Username, Password: r.Password}

	_, err = h.accountService.Create(c.Request.Context(), account)
	if err != nil {
		if errors.Is(err, apperrors.ErrorAccountAlreadyExists) {
			h.log.Warn("account already exists",
				slog.String(utils.Operation, op),
				slog.String("error", err.Error()))

			c.AbortWithStatusJSON(http.StatusConflict, ErrorResponse{Error: apperrors.ErrorAccountAlreadyExists.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *accountHandler) get(c *gin.Context) {
	const op = "api.get"
	l := h.log.With(slog.String(utils.Operation, op))

	aid := c.GetString("aid") //todo add special method for getting account id with error return

	acc, err := h.accountService.GetByID(c.Request.Context(), aid)
	if err != nil {
		if errors.Is(err, apperrors.ErrorAccountNotFound) {
			l.Warn("account not found",
				slog.String("account id", aid),
				slog.String("error", err.Error()))
			c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Error: apperrors.ErrorAccountNotFound.Error()})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, acc)
}

func (h *accountHandler) delete(c *gin.Context) {
	const op = "api.delete"
	l := h.log.With(slog.String(utils.Operation, op))
	aid := c.GetString("aid")

	l.Debug("deleting acc by id", slog.String("aid", aid))

	err := h.accountService.Delete(c.Request.Context(), aid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		l.Error("can't delete account",
			slog.String("error", err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "account was deleted",
	})
}
