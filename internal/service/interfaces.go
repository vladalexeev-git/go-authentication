package service

import (
	"context"
	"sso/internal/domain"
)

// Service
type Account interface {
	Create(ctx context.Context, acc domain.Account) (string, error)
	GetByID(ctx context.Context, aid string) (domain.Account, error)
	GetByEmail(ctx context.Context, email string) (domain.Account, error)
	Delete(ctx context.Context, aid string) error
}
type Auth interface {
	// EmailLogin creates new session using provided account email and password.
	//EmailLogin(ctx context.Context, email, password string, d Device) (domain.Session, error)
}

// Repository
type AccountRepo interface {
	Create(ctx context.Context, acc domain.Account) (string, error)
	FindByID(ctx context.Context, id string) (domain.Account, error)
	FindByEmail(ctx context.Context, email string) (domain.Account, error)
	Delete(ctx context.Context, id string) error
}

type SessionRepo interface {
	Create(ctx context.Context, session domain.Session) error
	FindByID(ctx context.Context, id string) (domain.Session, error)
	FindAll(ctx context.Context, aid string) ([]domain.Session, error)
}
