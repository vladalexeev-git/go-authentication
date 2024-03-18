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

	repo    AccountRepo
	session SessionRepo
}

func NewAccountService(cfg *config.Config, log *slog.Logger, repo AccountRepo, sess SessionRepo) *AccountService {

	return &AccountService{cfg: cfg, log: log, repo: repo, session: sess}
}

func (s *AccountService) Create(ctx context.Context, acc domain.Account) (string, error) {
	const op = "service.create"

	err := acc.GenPasswordHash()
	if err != nil {
		s.log.Error("can't gen password hash",
			slog.String(utils.Operation, op),
			slog.String("error", err.Error()))
		return "", fmt.Errorf("%s : %w", op, err)
	}

	aid, err := s.repo.Create(ctx, acc)
	if err != nil {
		s.log.Error("can't create account",
			slog.String(utils.Operation, op),
			slog.String("error", err.Error()))
		return "", fmt.Errorf("%s : %w", op, err)
	}
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
