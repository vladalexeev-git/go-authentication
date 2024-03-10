package service

import (
	"context"
	"fmt"
	"log/slog"
	"sso/config"
	"sso/internal/domain"
	"sso/pkg/utils"
)

type sessionService struct {
	log *slog.Logger
	cfg *config.Config

	repo SessionRepo
}

type Device struct {
	UserAgent string
	IP        string
}

func NewSessionService(log *slog.Logger, cfg *config.Config) *sessionService {
	return &sessionService{log: log, cfg: cfg}
}

func (s *sessionService) Create(ctx context.Context, aid, provider string, d Device) (domain.Session, error) {
	const op = "sessionservice.create"
	log := s.log.With(slog.String(utils.Operation, op))

	session, err := domain.NewSession(aid, provider, d.UserAgent, d.IP, s.cfg.Session.TTL)
	if err != nil {
		log.Error("can't create session",
			slog.String("error", err.Error()))
		return domain.Session{}, fmt.Errorf("%s: %w", op, err)
	}

	if err = s.repo.Create(ctx, session); err != nil {
		log.Error("can't create session",
			slog.String("error", err.Error()))
		return domain.Session{}, fmt.Errorf("%s: %w", op, err)
	}
	return session, nil
}
