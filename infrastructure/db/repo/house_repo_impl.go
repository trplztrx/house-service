package infrastructure

import (
	"context"
	"errors"
	"house-service/infrastructure/db/adapter"
	"house-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresHouseRepo struct {
	db           *pgxpool.Pool
	retryAdapter adapter.IPostgresRetryAdapter
}

func NewPostgresHouseRepo(db *pgxpool.Pool, retryAdapter adapter.IPostgresRetryAdapter) *PostgresHouseRepo {
	return &PostgresHouseRepo{
		db:           db,
		retryAdapter: retryAdapter,
	}
}

func (p *PostgresHouseRepo) Create(ctx context.Context, house *domain.House) (*domain.House, error) {
	var createdHouse domain.House
	query := `INSERT INTO houses(address, construction_date, developer, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING *`

	rows := p.retryAdapter.QueryRow(ctx, query,
		house.Address, house.Year,
		house.Developer, house.CreateAt,
		house.UpdateAt)
	defer rows.Close()

	err := rows.Scan(&createdHouse.HouseID, &createdHouse.Address,
		&createdHouse.Year, &createdHouse.Developer,
		&createdHouse.CreateAt, &createdHouse.UpdateAt)
	if err != nil {
		return &domain.House{}, err
	}

	return &createdHouse, nil
}


func (p *PostgresHouseRepo) GetByID(ctx context.Context, id int) (*domain.House, error) {
	var house domain.House
	query := `SELECT house_id, address, construction_date, developer, created_at, updated_at FROM houses WHERE house_id=$1`
	
	rows := p.retryAdapter.QueryRow(ctx, query, id)
	defer rows.Close()

	err := rows.Scan(&house.HouseID, &house.Address,
		&house.Year, &house.Developer,
		&house.CreateAt, &house.UpdateAt)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
			return &domain.House{}, nil
		}
		return &domain.House{}, err
	}
	
	return &house, nil
}

func (p *PostgresHouseRepo) Update(ctx context.Context, newHouseData *domain.House) error {
	query := `UPDATE houses SET address=$2, construction_date=$3, developer=$4, created_at=$5, updated_at=$6 WHERE house_id=$1`
	_, err := p.retryAdapter.Exec(ctx, query, newHouseData.HouseID, newHouseData.Address,
		newHouseData.Year, newHouseData.Developer,
		newHouseData.CreateAt, newHouseData.UpdateAt)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresHouseRepo) GetByAddress(ctx context.Context, address string) (*domain.House, error) {
	var house domain.House
	query := `SELECT house_id, address, construction_data, developer, created_at, updated_at 
			  FROM houses 
			  WHERE address = $1`

	rows := p.retryAdapter.QueryRow(ctx, query, address)
	defer rows.Close()

	err := rows.Scan(&house.HouseID, &house.Address,
		&house.Year, &house.Developer,
		&house.CreateAt, &house.UpdateAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &domain.House{}, nil // Если дом с таким адресом не найден, возвращаем nil и ошибку отсутствия данных
		}
		return &domain.House{}, err
	}

	return &house, nil
}

func (p *PostgresHouseRepo) GetApartmentsByID(ctx context.Context, id int, status string) ([]domain.Apartment, error) {
	var query string
	if status == domain.ModeratingStatus {
		query = `SELECT apartment_id, house_id, price, rooms, status FROM apartments WHERE house_id=$1 AND status=$2`
		} else {
		query = `SELECT aprtment_id, house_id, price, rooms, status FROM apartments WHERE house_id=$1`
	}

	rows, err := p.retryAdapter.Query(ctx, query, id, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apartments []domain.Apartment
	for rows.Next() {
		var apartment domain.Apartment
		err = rows.Scan(&apartment.ApartmentID, &apartment.HouseID, &apartment.Price, &apartment.Rooms, &apartment.Status)
		if err != nil {
			continue
		}
		apartments = append(apartments, apartment)
	}
	
	return apartments, err
}
func (p *PostgresHouseRepo) GetAll(ctx context.Context, offset int, limit int) ([]domain.House, error) {
	query := `SELECT house_id, address, construction_date, developer, created_at, updated_at FROM houses LIMIT $1 OFFSET $2`
	rows, err := p.retryAdapter.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var houses []domain.House
	for rows.Next() {
		var house domain.House
		err = rows.Scan(&house.HouseID, &house.Address, &house.Year,
			&house.Developer, &house.CreateAt, &house.UpdateAt)
		if err != nil {
			continue
		}
		houses = append(houses, house)
	}

	return houses, err
}

func (p *PostgresHouseRepo) DeleteByID(ctx context.Context, id int) error {
	query := `DELETE FROM houses WHERE id=$1`
	_, err := p.retryAdapter.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

