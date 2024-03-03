package service

import (
	"context"
	"fmt"
	"log/slog"
	"sso/config"
	"sso/internal/domain"
	"sso/pkg/utils"
)

type AccountService struct {
	cfg *config.Config
	log *slog.Logger

	repo AccountRepo
	//TODO: sessions

}

func NewAccountService(cfg *config.Config, log *slog.Logger) *AccountService {

	return &AccountService{cfg: cfg, log: log}
}

func (as *AccountService) Create(ctx context.Context, acc domain.Account) (string, error) {
	const op = "service.create"

	err := acc.GenPasswordHash()
	if err != nil {
		as.log.Error("can't gen password hash",
			slog.String(utils.Operation, op),
			slog.String("error", err.Error()))
		return "", fmt.Errorf("%s : %w", op, err)
	}

	aid, err := as.repo.Create(ctx, acc)
	if err != nil {
		as.log.Error("can't create account",
			slog.String(utils.Operation, op),
			slog.String("error", err.Error()))
		return "", fmt.Errorf("%s : %w", op, err)
	}
	return aid, nil
}

func (as *AccountService) GetByID(ctx context.Context, aid string) (domain.Account, error) {
	const op = "service.GetByID"

	acc, err := as.repo.FindByID(ctx, aid)
	if err != nil {
		as.log.Error("can't get account by id",
			slog.String(utils.Operation, op),
			slog.String("error", err.Error()))
		return domain.Account{}, fmt.Errorf("%s : %w", op, err)
	}

	return acc, nil
}

func (as *AccountService) GetByEmail(ctx context.Context, email string) (domain.Account, error) {
	const op = "service.FindByEmail"

	acc, err := as.repo.FindByEmail(ctx, email)
	if err != nil {
		as.log.Error("can't get account by email",
			slog.String(utils.Operation, op),
			slog.String("error", err.Error()))
		return domain.Account{}, fmt.Errorf("%s : %w", op, err)
	}

	return acc, nil
}

//Delete