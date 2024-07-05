package service

import (
	"context"
	"fmt"
	"go-authentication/config"
	"go-authentication/internal/domain"
	"go-authentication/pkg/utils"
	"log/slog"
)

type AccountService struct {
	cfg *config.Config
	log *slog.Logger

	repo    AccountRepo
	session SessionRepo
}

func NewAccountService(cfg *config.Config, log *slog.Logger, repo AccountRepo, sess SessionRepo) *AccountService {

	return &AccountService{cfg: cfg, log: log, repo: repo, session: sess}
}

func (s *AccountService) Create(ctx context.Context, acc domain.Account) (string, error) {
	const op = "service.create"
	l := s.log.With(slog.String(utils.Operation, op))

	err := acc.GenPasswordHash()
	if err != nil {
		l.Error("can't gen password hash",
			slog.String("error", err.Error()))
		return "", fmt.Errorf("%s : %w", op, err)
	}

	aid, err := s.repo.Create(ctx, acc)
	if err != nil {
		return "", fmt.Errorf("%s : %w", op, err)
	}

	l.Info("account created successfully", slog.String("account_id", aid))

	return aid, nil
}

func (s *AccountService) GetByID(ctx context.Context, aid string) (domain.Account, error) {
	const op = "service.GetByID"

	acc, err := s.repo.FindByID(ctx, aid)
	if err != nil {
		return domain.Account{}, fmt.Errorf("%s : %w", op, err)
	}

	return acc, nil
}

func (s *AccountService) GetByEmail(ctx context.Context, email string) (domain.Account, error) {
	const op = "service.FindByEmail"

	acc, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return domain.Account{}, fmt.Errorf("%s : %w", op, err)
	}

	return acc, nil
}

func (s *AccountService) Delete(ctx context.Context, aid string) error {
	const op = "service.Delete"

	err := s.repo.Delete(ctx, aid)
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}
	return nil
}
