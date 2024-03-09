package v1

import (
	"context"
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

func NewAccount(handler *gin.RouterGroup, log *slog.Logger, cfg *config.Config, accService service.Account) {
	h := &accountHandler{log: log, cfg: cfg, accountService: accService}

	handler.POST("/create", h.create)
	handler.GET("/:id", h.get)
}

//TODO: Create special errors understandable for users
//TODO: Think, maybe import domain models in this layer is a bad practice?

//TODO: add context instead of context.TODO

func (ah *accountHandler) create(c *gin.Context) {
	const op = "api.create"
	var r accountCreateRequest

	err := c.BindJSON(&r)
	if err != nil {
		ah.log.Error("can't unmarshal account request",
			slog.String(utils.Operation, op),
			slog.String("error", err.Error()))

		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: apperrors.ErrorValidate.Error()})
		return
	}

	account := domain.Account{Email: r.Email, Username: r.Username, Password: r.Password}

	_, err = ah.accountService.Create(context.TODO(), account)
	if err != nil {
		if errors.Is(err, apperrors.ErrorAccountAlreadyExists) {
			ah.log.Error("account already exists",
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

func (ah *accountHandler) get(c *gin.Context) {
	const op = "api.get"
	aid := c.Param("id")

	acc, err := ah.accountService.GetByID(context.TODO(), aid)
	if err != nil {
		if errors.Is(err, apperrors.ErrorAccountNotFound) {
			ah.log.Error("account not found",
				slog.String(utils.Operation, op),
				slog.String("error", err.Error()))
			c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Error: apperrors.ErrorAccountNotFound.Error()})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, acc)
}
