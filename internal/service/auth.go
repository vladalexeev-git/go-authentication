package service

import (
	"github.com/golang-jwt/jwt"
	"log/slog"
)

type authService struct {
	log     *slog.Logger
	token   jwt.Token
	account Account
}

func NewAuthService(log *slog.Logger, token jwt.Token, account Account) *authService {
	return &authService{log: log, token: token, account: account}
}

//func (j *authService) EmailLogin(ctx context.Context, email, password string, d Device) (domain.Session, error) {
