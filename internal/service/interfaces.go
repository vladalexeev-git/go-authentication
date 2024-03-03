package service

import (
	"context"
	"sso/internal/domain"
)

type Account interface {
	Create(ctx context.Context, acc domain.Account) (string, error)
	GetByID(ctx context.Context, aid string) (domain.Account, error)
	GetByEmail(ctx context.Context, email string) (domain.Account, error)
}

type AccountRepo interface {
	Create(ctx context.Context, acc domain.Account) (string, error)
	FindByID(ctx context.Context, id string) (domain.Account, error)
	FindByEmail(ctx context.Context, email string) (domain.Account, error)
	Delete(ctx context.Context, id string) error

	//TODO: Update, Delete
}
