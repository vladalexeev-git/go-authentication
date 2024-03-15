package service

import (
	"context"
	"github.com/golang-jwt/jwt"
	"log/slog"
	"sso/internal/domain"
	"sso/pkg/utils"
)

type authService struct {
	log     *slog.Logger
	token   jwt.Token
	account Account
	session sessionService
}

func NewAuthService(log *slog.Logger, token jwt.Token, account Account) *authService {
	return &authService{log: log, token: token, account: account}
}

func (as *authService) EmailLogin(ctx context.Context, email, password string, d Device) (domain.Session, error) {
	const op = "EmailLogin"
	var s domain.Session
	l := as.log.With(slog.String(utils.Operation, op))

	acc, err := as.account.GetByEmail(ctx, email)
	if err != nil {
		//todo return appropriate error and log
	}

	acc.Password = password
	err = acc.CompareHashAndPassword()
	if err != nil {
		l.Error("password is incorrect")
		return s, err
	}

	return s, nil
}
