package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/reecevinto/coaches-revenue-intelligences-saas/internal/users"
)

// UserRepository implements users.Repository using PostgreSQL
type UserRepository struct {
	pool *pgxpool.Pool
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

// Create inserts a new user into the database
func (r *UserRepository) Create(ctx context.Context, user *users.User) error {
	query := `
		INSERT INTO users (id, account_id, email, password, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.pool.Exec(
		ctx,
		query,
		user.ID,
		user.AccountID,
		user.Email,
		user.Password,
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
	)
	return err
}

// GetByEmail fetches a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*users.User, error) {
	query := `
		SELECT id, account_id, email, password, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	row := r.pool.QueryRow(ctx, query, email)
	var user users.User
	err := row.Scan(
		&user.ID,
		&user.AccountID,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ✅ GetByID fetches a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*users.User, error) {
	query := `
		SELECT id, account_id, email, password, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	var user users.User
	err := row.Scan(
		&user.ID,
		&user.AccountID,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
