package service

import (
	"context"
	"fmt"
	"go-authentication/config"
	"go-authentication/internal/apperrors"
	"go-authentication/internal/domain"
)

type sessionService struct {
	cfg *config.Config

	repo SessionRepo
}

type Device struct {
	UserAgent string
	IP        string
}

func NewSessionService(cfg *config.Config, repo SessionRepo) *sessionService {
	return &sessionService{cfg: cfg, repo: repo}
}

func (s *sessionService) Create(ctx context.Context, aid, provider string, d Device) (domain.Session, error) {
	const op = "sessionservice.create"

	session, err := domain.NewSession(aid, provider, d.UserAgent, d.IP, s.cfg.Session.TTL)
	if err != nil {

		return domain.Session{}, fmt.Errorf("%s: %w", op, err)
	}

	if err = s.repo.Create(ctx, session); err != nil {
		return domain.Session{}, fmt.Errorf("%s: %w", op, err)
	}
	return session, nil
}

func (s *sessionService) Get(ctx context.Context, sid string) (domain.Session, error) {
	const op = "sessionservice.get"

	session, err := s.repo.FindByID(ctx, sid)
	if err != nil {
		return domain.Session{}, fmt.Errorf("%s: %w", op, err)
	}

	return session, nil
}

func (s *sessionService) GetAll(ctx context.Context, aid string) ([]domain.Session, error) {
	const op = "sessionservice.getall"

	sessions, err := s.repo.FindAll(ctx, aid)
	if err != nil {
		return []domain.Session{}, fmt.Errorf("%s: %w", op, err)
	}

	return sessions, nil
}

func (s *sessionService) Terminate(ctx context.Context, curSid string, reqSid string) error {
	const op = "sessionservice.terminate"

	if curSid == reqSid {
		return fmt.Errorf("%s: %w", op, apperrors.ErrorCurrentSessionTerminating)
	}

	if err := s.repo.Delete(ctx, reqSid); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *sessionService) TerminateAll(ctx context.Context, aid, sid string) error {
	const op = "sessionservice.terminateAll"

	if err := s.repo.DeleteAll(ctx, aid, sid); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
