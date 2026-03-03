package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/reecevinto/coaches-revenue-intelligences-saas/internal/accounts"
)

type AccountRepository struct {
	pool *pgxpool.Pool
}

func NewAccountRepository(pool *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{pool: pool}
}

func (r *AccountRepository) Create(ctx context.Context, account *accounts.Account) error {
	query := `
		INSERT INTO accounts (id, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.pool.Exec(
		ctx,
		query,
		account.ID,
		account.Name,
		account.CreatedAt,
		account.UpdatedAt,
	)

	return err
}

func (r *AccountRepository) GetByID(ctx context.Context, id string) (*accounts.Account, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM accounts
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)

	var account accounts.Account

	err := row.Scan(
		&account.ID,
		&account.Name,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &account, nil
}
