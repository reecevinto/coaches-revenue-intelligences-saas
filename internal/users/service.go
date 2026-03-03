package users

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(ctx context.Context, accountID, email, password, role string) (*User, error) {
	user := &User{
		ID:        uuid.NewString(),
		AccountID: accountID,
		Email:     email,
		Password:  password, // will hash later
		Role:      role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
