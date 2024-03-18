package service

import (
	"context"
	"fmt"
	"log/slog"
	"sso/internal/domain"
	"sso/pkg/JWT"
	"sso/pkg/utils"
)

type authService struct {
	log   *slog.Logger
	token JWT.Token

	account Account
	session Session
}

func NewAuthService(log *slog.Logger, token JWT.Token, account Account, session Session) *authService {
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
		slog.String("email", email),
		slog.String("id", a.ID),
		slog.String("password hash", a.PasswordHash))

	a.Password = password
	err = a.CompareHashAndPassword()
	if err != nil {
		l.Error("can't login", slog.String("error", err.Error()))
		l.Debug("", slog.String("hashed password", a.PasswordHash))
		return domain.Session{}, fmt.Errorf("%s: %w", op, err)
	}

	l.Debug("session is ready to creating", slog.Any("account", a))
	//creating a session
	sess, err := s.session.Create(ctx, a.ID, a.Email, d)
	if err != nil {
		return domain.Session{}, fmt.Errorf("%s: %w", op, err)
	}

	return sess, nil
}
