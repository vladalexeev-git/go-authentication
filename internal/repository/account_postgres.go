package repository

import (
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
	"sso/config"
	"sso/internal/domain"
)

type accountRepo struct {
	log *slog.Logger
	cfg *config.Config

	//TODO: gormClient
	db *gorm.DB
}

func NewAccountRepo(log *slog.Logger, cfg *config.Config, db *gorm.DB) *accountRepo {
	return &accountRepo{
		log: log,
		cfg: cfg,
		db:  db,
	}
}

func (ar *accountRepo) CreateAccount(acc domain.Account) (string, error) {

	return "", nil
}
