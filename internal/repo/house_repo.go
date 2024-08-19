package repo

import (
	"context"
	"house-service/internal/domain"
)

type HouseRepo interface {
	Create(ctx context.Context, house *domain.House) (*domain.House, error)
	GetByID(ctx context.Context, id int) (*domain.House, error)
	GetByAddress(ctx context.Context, address string) (*domain.House, error)
	GetApartmentsByID(ctx context.Context, id int, status string) ([]domain.Apartment, error)
	GetAll(ctx context.Context, offset int, limit int) ([]domain.House, error)
	DeleteByID(ctx context.Context, id int) error
	Update(ctx context.Context, newHouse *domain.House) error
}
