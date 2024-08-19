package repo

import (
	"context"
	"house-service/internal/domain"

	"github.com/google/uuid"
)

type UserRepo interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetAll(ctx context.Context, offset int, limit int) ([]domain.User, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	// Update(ctx context.Context, newUser *domain.User) (*domain.User, error)
}
