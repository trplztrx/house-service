package repo

import (
	"context"
	"house-service/internal/domain"

	"github.com/google/uuid"
)

// реализация на урвоне infrastructure

type ApartmentRepo interface {
	Create(ctx context.Context, apartment *domain.Apartment) (*domain.Apartment, error)
	GetByID(ctx context.Context, id int, houseID int) (*domain.Apartment, error)
	DeleteByID(ctx context.Context, apartmentID int, houseID int) error
	Update(ctx context.Context, moderatorID uuid.UUID, newApartment *domain.Apartment) (*domain.Apartment, error)
	GetAll(ctx context.Context, offset int, limit int) ([]domain.Apartment, error)
}