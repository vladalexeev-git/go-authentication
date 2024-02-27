package v1

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"sso/config"
	"sso/internal/domain"
	"sso/internal/service"
	"sso/pkg/utils"
)

type accountHandler struct {
	//TODO: interface accountService
	log *slog.Logger
	cfg *config.Config

	accountService service.AccountService
}

func NewAccount(handler *gin.RouterGroup, log *slog.Logger, cfg *config.Config) {
	h := &accountHandler{log: log, cfg: cfg}

	handler.POST("/", h.create)
}

// request model
type accountCreateRequest struct {
	Email    string `json:"email" binding:"required,email,lte=255"`
	Username string `json:"username" binding:"required,alphanum,gte=4,lte=16"`
	Password string `json:"password" binding:"required,gte=8,lte=64"`
}

//TODO: Create special errors understandable for users
//TODO: Think, maybe import domain models in this layer is a bad practice?

func (ah *accountHandler) create(c *gin.Context) {
	const op = "api/CreateAccount"
	var r accountCreateRequest

	err := c.BindJSON(&r)
	if err != nil {
		ah.log.Error("can't unmarshal account request", slog.String(utils.Operation, op), slog.String("error", err.Error()))

		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	account := domain.Account{Email: r.Email, Username: r.Username, Password: r.Password}

	err = ah.accountService.Create(c.Request.Context(), account)
	if err != nil {

	}
}
