package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go-authentication/internal/apperrors"
	"go-authentication/internal/domain"
	"go-authentication/pkg/postgres"
	"go-authentication/pkg/utils"
	"log/slog"
)

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
func (r *accountRepo) Create(ctx context.Context, acc domain.Account) (string, error) {
	const op = "repository.accountRepo.Create"
	l := r.log.With(slog.String(utils.Operation, op))

	sql, args, err := r.pg.Builder.
		Insert(_accTable).
		Columns("username, email, password").
		Values(acc.Username, acc.Email, acc.PasswordHash).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		l.Error("pg.builder: bad insert query",
			slog.String("error", err.Error()))

		return "", fmt.Errorf("%s : %w", op, err)
	}

	//aid, err := ar.db.Pool.Exec(ctx, sql, args...)
	var aid string

	if err = r.pg.Pool.QueryRow(ctx, sql, args...).Scan(&aid); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				l.Error("queryrow uniq violation", slog.String("error", err.Error()))
				return "", fmt.Errorf("%s: %w", op, apperrors.ErrorAccountAlreadyExists)
			}
			l.Error("queryrow error", slog.String("error", err.Error()))
			return "", fmt.Errorf("%s : %w", op, err)
		}
	}
	return aid, nil
}

// FindByID ...
func (r *accountRepo) FindByID(ctx context.Context, aid string) (domain.Account, error) {
	const op = "repository.accountRepo.GetByID"
	l := r.log.With(slog.String(utils.Operation, op))

	sql, args, err := r.pg.Builder.
		Select("username", "email", "password", "created_at", "updated_at").
		From(_accTable).
		Where(squirrel.Eq{"id": aid}).
		ToSql()
	if err != nil {
		l.Error("builder - bad select by id query",
			slog.Any("args", args),
			slog.String("sql", sql),
			slog.String("error", err.Error()))
		return domain.Account{}, fmt.Errorf("%s : %w", op, err)
	}

	var acc = domain.Account{ID: aid}

	if err = r.pg.Pool.QueryRow(ctx, sql, args...).Scan(
		&acc.Username,
		&acc.Email,
		&acc.PasswordHash,
		&acc.CreatedAt,
		&acc.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			l.Error("account not found",
				slog.String("error", err.Error()))
			return domain.Account{}, fmt.Errorf("%s: %w", op, apperrors.ErrorAccountNotFound)
		}
		l.Error("bad queryRow or scan",
			slog.String("error", err.Error()))
		return domain.Account{}, fmt.Errorf("%s : %w", op, err)
	}

	return acc, nil
}

// FindByEmail ...
func (r *accountRepo) FindByEmail(ctx context.Context, email string) (domain.Account, error) {
	const op = "repository.accountRepo.GetByEmail"
	l := r.log.With(slog.String(utils.Operation, op))

	sql, args, err := r.pg.Builder.
		Select("id", "username", "password", "created_at", "updated_at").
		From(_accTable).
		Where(squirrel.Eq{"email": email}).
		ToSql()
	if err != nil {
		l.Error("builder - bad select by email query",
			slog.Any("args", args),
			slog.String("sql", sql),
			slog.String("error", err.Error()))
		return domain.Account{}, fmt.Errorf("%s : %w", op, err)
	}

	var acc = domain.Account{
		Email: email,
	}

	if err = r.pg.Pool.QueryRow(ctx, sql, args...).Scan(
		&acc.ID,
		&acc.Username,
		&acc.PasswordHash,
		&acc.CreatedAt,
		&acc.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			l.Error("account not found",
				slog.String("error", err.Error()))
			return domain.Account{}, fmt.Errorf("%s: %w", op, apperrors.ErrorAccountNotFound)
		}
		l.Error("bad queryRow or scan",
			slog.String("error", err.Error()))
		return domain.Account{}, fmt.Errorf("%s : %w", op, err)
	}
	return acc, nil
}

// Delete ...
func (r *accountRepo) Delete(ctx context.Context, aid string) error {
	const op = "repository.accountRepo.Delete"
	l := r.log.With(slog.String(utils.Operation, op))

	sql, args, err := r.pg.Builder.
		Delete(_accTable).
		Where(squirrel.Eq{"id": aid}).
		ToSql()

	if err != nil {
		l.Error("builder - bad delete query",
			slog.String("sql", sql),
			slog.Any("args", args),
			slog.String("error", err.Error()))
		return fmt.Errorf("%s : %w", op, err)
	}

	ct, err := r.pg.Pool.Exec(ctx, sql, args...)
	r.log.Debug("returned result",
		slog.Int64("count", ct.RowsAffected()),
		slog.String("string", ct.String()))
	if err != nil {
		l.Error("pool.exec", slog.String("error", err.Error()))
		return fmt.Errorf("%s : %w", op, err)
	}
	return err
}
