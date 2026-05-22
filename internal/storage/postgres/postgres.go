package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lfssxxx/auth_service/internal/domain/models"
	"github.com/lfssxxx/auth_service/internal/storage"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Storage {
	return &Storage{
		pool: pool,
	}
}
func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.postgres.SaveUser"

	query := `
	INSERT INTO users(email, pass_hash)
	VALUES ($1, $2)
	RETURNING id;`
	row := s.pool.QueryRow(ctx, query, email, passHash)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("scan error at: %s: %w", op, storage.ErrUserNotFound)
		}
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return 0, fmt.Errorf("user already exists: %s: %w", op, storage.ErrUserExists)
			}
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.postgres.User"

	query := `
	SELECT id, email, pass_hash
	FROM users
	WHERE email=$1;
	`

	row := s.pool.QueryRow(ctx, query, email)

	var userModel models.User

	err := row.Scan(
		&userModel.ID,
		&userModel.Email,
		&userModel.PassHash,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return userModel, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.postgres.IsAdmin"

	query := `
	SELECT is_admin
	FROM users
	WHERE id=$1;
	`

	row := s.pool.QueryRow(ctx, query, userID)

	var isAdmin bool

	err := row.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil

}

func (s *Storage) App(ctx context.Context, appID int) (models.App, error) {
	const op = "storage.postgres.App"

	query := `
	SELECT id, name, secret FROM apps WHERE id=$1;
	`

	row := s.pool.QueryRow(ctx, query, appID)

	var app models.App

	err := row.Scan(
		&app.ID,
		&app.Name,
		&app.Secret,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}

		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}
