package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"house-service/infrastructure/db/adapter"
	"house-service/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepo struct {
	db           *pgxpool.Pool
	retryAdapter adapter.IPostgresRetryAdapter
}

func NewPostgresUserRepo(db *pgxpool.Pool, retryAdapter adapter.IPostgresRetryAdapter) *PostgresUserRepo {
	return &PostgresUserRepo{
		db:           db,
		retryAdapter: retryAdapter,
	}
}

func (p *PostgresUserRepo) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users(user_id, mail, password, role) VALUES ($1, $2, $3, $4)`
	_, err := p.retryAdapter.Exec(ctx, query, user.UserID, user.Email, user.Password, user.Type)
	if err != nil {
		return err
	}
	return nil
}


// func (p *PostgresUserRepo) Update(ctx context.Context, newUserData *domain.User) (*domain.User, error) {
// 	query := `UPDATE users SET mail=$2, password=$3, role=$4 WHERE user_id=$1`
// 	user, err := p.retryAdapter.Exec(ctx, query, newUserData.UserID, newUserData.Email, newUserData.Password, newUserData.Type)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &user
// }

func (p *PostgresUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	query := `SELECT user_id, mail, password, role FROM users WHERE user_id=$1`
	rows := p.retryAdapter.QueryRow(ctx, query, id)
	defer rows.Close()
	err := rows.Scan(&user.UserID, &user.Email, &user.Password, &user.Type)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &domain.User{}, nil 
		}
		return &domain.User{}, err
	}
	return &user, nil
}
func (p *PostgresUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	query := `SELECT user_id, mail, password, role FROM users WHERE mail = $1`
	rows := p.retryAdapter.QueryRow(ctx, query, email)
	defer rows.Close()

	err := rows.Scan(&user.UserID, &user.Email, &user.Password, &user.Type)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (p *PostgresUserRepo) GetAll(ctx context.Context, offset int, limit int) ([]domain.User, error) {
	query := `SELECT user_id, mail, password, role FROM users LIMIT $1 OFFSET $2`
	rows, err := p.retryAdapter.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("user repo: get all error: %v", err)
	}
	defer rows.Close()
	
	var users []domain.User
	for rows.Next() {
		var user domain.User
		err = rows.Scan(&user.UserID, &user.Email, &user.Password, &user.Type)
		if err != nil {
			continue
		}
		users = append(users, user)
	}
	return users, err
}

func (p *PostgresUserRepo) DeleteByID(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE user_id=$1`
	_, err := p.retryAdapter.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

