package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"house-service/infrastructure/db/adapter"
	"house-service/internal/domain"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresapartmentRepo struct {
	db           *pgxpool.Pool
	retryAdapter adapter.IPostgresRetryAdapter
}

func NewPostgresApartmentRepo(db *pgxpool.Pool, retryAdapter adapter.IPostgresRetryAdapter) *PostgresapartmentRepo {
	return &PostgresapartmentRepo{
		db:           db,
		retryAdapter: retryAdapter,
	}
}

func (p *PostgresapartmentRepo) Create(ctx context.Context, apartment *domain.Apartment) (*domain.Apartment, error) {
	var createdApartment domain.Apartment

	tx, err := p.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return &domain.Apartment{}, fmt.Errorf("postgres apartment repo: create error: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	query := `INSERT INTO apartments(apartment_id, house_id, user_id, price, rooms, status)
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING apartment_id, house_id, user_id, price, rooms, status`
	err = tx.QueryRow(ctx, query, apartment.ApartmentID, apartment.HouseID, apartment.OwnerID,
		apartment.Price, apartment.Rooms, domain.CreatedStatus).Scan(&createdApartment.ApartmentID,
		&createdApartment.HouseID, &createdApartment.OwnerID, &createdApartment.Price, &createdApartment.Rooms,
		&createdApartment.Status)
	if err != nil {
		tx.Rollback(ctx)
		return &domain.Apartment{}, fmt.Errorf("postgres apartment repo: create error: %w", err)
	}

	date := time.Now()
	query = `UPDATE houses SET update_apartment_date=$1 WHERE house_id=$2`
	_, err = tx.Exec(ctx, query, date, createdApartment.HouseID)
	if err != nil {
		tx.Rollback(ctx)
		return &domain.Apartment{}, fmt.Errorf("postgres apartment repo: create error: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return &domain.Apartment{}, fmt.Errorf("postgres apartment repo: create error: %w", err)
	}

	return &createdApartment, nil
}

func (p *PostgresapartmentRepo) DeleteByID(ctx context.Context, apartmentID int, houseID int) error {
	query := `DELETE FROM apartments WHERE apartment_id=$1 AND house_id=$2`
	_, err := p.retryAdapter.Exec(ctx, query, apartmentID, houseID)
	if err != nil {
		return fmt.Errorf("postgres apartment repo: delete by id error: %w", err)
	}
	return nil
}

func (p *PostgresapartmentRepo) Update(ctx context.Context, moderatorID uuid.UUID, newApartmentData *domain.Apartment) (*domain.Apartment, error) {
	var apartment domain.Apartment

	query := `SELECT apartment_id, house_id, user_id, price, rooms, status FROM update_status($1, $2, $3, $4)`
	rows := p.retryAdapter.QueryRow(ctx, query, newApartmentData.Status,
		newApartmentData.ApartmentID, newApartmentData.HouseID, moderatorID)
	defer rows.Close()

	err := rows.Scan(&apartment.ApartmentID, &apartment.HouseID, &apartment.OwnerID,
		&apartment.Price, &apartment.Rooms, &apartment.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &domain.Apartment{}, nil
		}
		return &domain.Apartment{}, fmt.Errorf("postgres apartment repo: update error: %w", err)
	}

	return &apartment, nil
}

func (p *PostgresapartmentRepo) GetByID(ctx context.Context, apartmentID int, houseID int) (*domain.Apartment, error) {
	var apartment domain.Apartment

	query := `SELECT apartment_id, house_id, user_id, price, rooms, status
	FROM apartments WHERE apartment_id=$1 AND house_id=$2`
	rows := p.retryAdapter.QueryRow(ctx, query, apartmentID, houseID)
	defer rows.Close()

	err := rows.Scan(&apartment.ApartmentID, &apartment.HouseID, &apartment.OwnerID,
		&apartment.Price, &apartment.Rooms, &apartment.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &domain.Apartment{}, nil
		}
		return &domain.Apartment{}, fmt.Errorf("postgres apartment repo: get by id error: %w", err)
	}

	return &apartment, nil
}

func (p *PostgresapartmentRepo) GetAll(ctx context.Context, offset int, limit int) ([]domain.Apartment, error) {
	query := `SELECT apartment_id, house_id, user_id, price, rooms, status FROM apartments LIMIT $1 OFFSET $2`
	rows, err := p.retryAdapter.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("postgres apartment repo: get all error: %w", err)
	}
	defer rows.Close()

	var apartments []domain.Apartment
	for rows.Next() {
		var apartment domain.Apartment
		err = rows.Scan(&apartment.ApartmentID, &apartment.HouseID, &apartment.OwnerID,
			&apartment.Price, &apartment.Rooms, &apartment.Status)
		if err != nil {
			continue
		}
		apartments = append(apartments, apartment)
	}

	return apartments, err
}
