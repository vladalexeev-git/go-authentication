package service

import (
	"context"
	"sso/internal/domain"
)

// Services

type Account interface {
	Create(ctx context.Context, acc domain.Account) (string, error)
	GetByID(ctx context.Context, aid string) (domain.Account, error)
	GetByEmail(ctx context.Context, email string) (domain.Account, error)
	Delete(ctx context.Context, aid string) error
}

type Session interface {
	Create(ctx context.Context, aid, provider string, d Device) (domain.Session, error)
	Get(ctx context.Context, sid string) (domain.Session, error)
	Terminate(ctx context.Context, sid string) error
}

type Auth interface {
	// EmailLogin creates new session using provided account email and password.
	EmailLogin(ctx context.Context, email, password string, d Device) (domain.Session, error)
	Logout(ctx context.Context, sid string) error
}

// Repositories:

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
	Delete(ctx context.Context, sid string) error
	DeleteAll(ctx context.Context, aid, currSid string) error
}
