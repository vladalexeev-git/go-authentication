package service

import (
	"context"
	"fmt"
	"log/slog"
	"sso/config"
	"sso/internal/domain"
	"sso/pkg/utils"
)

//TODO: Is it a bad practice to put logger into service layer?

type AccountService struct {
	cfg *config.Config
	log *slog.Logger
	//TODO: repo repositoryAccount
	//TODO: sessions

	//repository
}

func NewAccountService(cfg *config.Config, log *slog.Logger) *AccountService {

	return &AccountService{cfg: cfg, log: log}
}

func (as *AccountService) Create(ctx context.Context, acc domain.Account) error {
	const op = "service/Create"

	err := acc.GenPasswordHash()
	if err != nil {
		as.log.Error("can't gen password hash", slog.String(utils.Operation, op), slog.String("error", err.Error()))
		return fmt.Errorf("%s : %w", op, err)
	}
	return nil
}
