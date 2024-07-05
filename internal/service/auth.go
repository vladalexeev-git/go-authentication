package service

import (
	"context"
	"fmt"
	"go-authentication/internal/domain"
	"go-authentication/pkg/utils"
	"log/slog"
)

type authService struct {
	log     *slog.Logger
	token   Token
	account Account
	session Session
}

func NewAuthService(log *slog.Logger, token Token, account Account, session Session) *authService {
	return &authService{log: log, token: token, account: account, session: session}
}

func (s *authService) EmailLogin(ctx context.Context, email, password string, d Device) (domain.Session, error) {
	const op = "auth.emailLogin"
	l := s.log.With(slog.String(utils.Operation, op))

	//fetching the account
	a, err := s.account.GetByEmail(ctx, email)
	if err != nil {
		return domain.Session{}, fmt.Errorf("%s: %w", op, err)
	}
	l.Debug("account found",
		slog.Any("account", a))

	a.Password = password
	err = a.CompareHashAndPassword()
	if err != nil {
		l.Error("can't login", slog.String("error", err.Error()))
		l.Debug("", slog.String("hashed password", a.PasswordHash))
		return domain.Session{}, fmt.Errorf("%s: %w", op, err)
	}

	//creating a session
	sess, err := s.session.Create(ctx, a.ID, a.Email, d)
	if err != nil {
		return domain.Session{}, fmt.Errorf("%s: %w", op, err)
	}

	return sess, nil
}

func (s *authService) Logout(ctx context.Context, sid string) error {
	const op = "auth.logout"

	err := s.session.Terminate(ctx, sid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *authService) NewAccessToken(ctx context.Context, sub, password string) (string, error) {
	const op = "auth.AccessToken"

	a, err := s.account.GetByID(ctx, sub)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	a.Password = password
	err = a.CompareHashAndPassword()
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	t, err := s.token.New(sub)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return t, nil
}

func (s *authService) ParseAccessToken(ctx context.Context, token string) (string, error) {
	const op = "auth.ParseAccessToken"

	aid, err := s.token.Parse(token)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return aid, nil
}
