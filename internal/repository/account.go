package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"log/slog"
	"sso/internal/domain"
	"sso/pkg/apperrors"
	"sso/pkg/postgres"
	"sso/pkg/utils"
)

// TODO: think about delete method may be should delete by whole acc or email too
const _accTable = "accounts"

type accountRepo struct {
	log *slog.Logger
	pg  *postgres.Postgres
}

func NewAccountRepo(log *slog.Logger, db *postgres.Postgres) *accountRepo {
	return &accountRepo{
		log: log,
		pg:  db,
	}
}

// Create ...
func (ar *accountRepo) Create(ctx context.Context, acc domain.Account) (string, error) {
	const op = "repository.accountRepo.Create"

	sql, args, err := ar.pg.Builder.
		Insert("username, email, password").
		Values(acc.Username, acc.Email, acc.Password).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		ar.log.Error("builder - bad insert query",
			slog.String(utils.Operation, op),
			slog.String("error", err.Error()))

		return "", fmt.Errorf("%s : %w", op, err)
	}

	//aid, err := ar.db.Pool.Exec(ctx, sql, args...)
	var aid string

	if err = ar.pg.Pool.QueryRow(ctx, sql, args...).Scan(&aid); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return "", fmt.Errorf("%s: %w", op, apperrors.ErrorAccountAlreadyExists)
			}
			return "", fmt.Errorf("%s : %w", op, err)
		}
	}
	return aid, nil
}

// FindByID ...
func (ar *accountRepo) FindByID(ctx context.Context, aid string) (domain.Account, error) {
	const op = "repository.accountRepo.GetByID"

	sql, args, err := ar.pg.Builder.
		Select("username", "email", "password", "created_at", "updated_at").
		From(_accTable).
		Where(squirrel.Eq{"id": aid}).
		ToSql()
	if err != nil {
		ar.log.Error("builder - bad select by id query",
			slog.String(utils.Operation, op),
			slog.String("error", err.Error()))
		return domain.Account{}, fmt.Errorf("%s : %w", op, err)
	}

	var acc = domain.Account{ID: aid}

	if err = ar.pg.Pool.QueryRow(ctx, sql, args...).Scan(
		&acc.Username,
		&acc.Email,
		&acc.Password,
		&acc.CreatedAt,
		&acc.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Account{}, fmt.Errorf("%s: %w", op, apperrors.ErrorAccountNotFound)
		}
		return domain.Account{}, fmt.Errorf("%s : %w", op, err)
	}

	return acc, nil
}

// FindByEmail ...
func (ar *accountRepo) FindByEmail(ctx context.Context, email string) (domain.Account, error) {
	const op = "repository.accountRepo.GetByEmail"

	sql, args, err := ar.pg.Builder.
		Select("username", "email", "password", "created_at", "updated_at").
		From(_accTable).
		Where(squirrel.Eq{"email": email}).
		ToSql()
	if err != nil {
		ar.log.Error("builder - bad select by email query",
			slog.String(utils.Operation, op),
			slog.String("error", err.Error()))
		return domain.Account{}, fmt.Errorf("%s : %w", op, err)
	}

	var acc = domain.Account{
		Email: email,
	}

	if err = ar.pg.Pool.QueryRow(ctx, sql, args...).Scan(
		&acc.Username,
		&acc.Email,
		&acc.Password,
		&acc.CreatedAt,
		&acc.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Account{}, fmt.Errorf("%s: %w", op, apperrors.ErrorAccountNotFound)
		}
		return domain.Account{}, fmt.Errorf("%s : %w", op, err)
	}
	return acc, nil
}

// Delete ...
func (ar *accountRepo) Delete(ctx context.Context, aid string) error {
	const op = "repository.accountRepo.Delete"

	sql, args, err := ar.pg.Builder.
		Delete(_accTable).
		Where(squirrel.Eq{"id": aid}).
		ToSql()

	if err != nil {
		ar.log.Error("builder - bad delete query",
			slog.String(utils.Operation, op),
			slog.String("error", err.Error()))
		return fmt.Errorf("%s : %w", op, err)
	}

	_, err = ar.pg.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}
	return err
}
