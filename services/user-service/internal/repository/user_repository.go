package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sergssclaude/health-tracker/user-service/internal/model"
)

var ErrUserNotFound = errors.New("User not found")
var ErrEmailAlreadyExists = errors.New("Email already exists")

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, userId int) (*model.User, error)
	GetByMail(ctx context.Context, email string) (*model.User, error)
}

type postgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &postgresUserRepository{pool: pool}
}

func (r *postgresUserRepository) Create(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users (email, password_hash)
	VALUES ($1, $2) RETURNING id`
	err := r.pool.QueryRow(ctx, query, user.Email, user.PasswordHash).Scan(&user.Id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrEmailAlreadyExists
		}
	}
	return nil
}

func (r *postgresUserRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	query := `SELECT id, email, password_hash, name, created_at, updated_at
	FROM users WHERE id = $1`
	var user model.User
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.Id,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *postgresUserRepository) GetByMail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, email, password_hash, name, created_at, updated_at
	FROM users WHERE email = $1`
	var user model.User
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.Id,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}
