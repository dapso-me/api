package customer_postgresql_repository

import (
	"api/internal/common"
	"api/internal/domain/customer"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repo struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) customer.Repository {
	return &repo{
		db: db,
	}
}

func (r *repo) FindOneByID(c context.Context, customerID uuid.UUID) (*customer.Customer, error) {
	var customerEntity customer.Customer

	sql := "SELECT id, username, password, name, role, created_at FROM customers.customers WHERE id = $1"
	row := r.db.QueryRow(c, sql, customerID)

	err := row.Scan(
		&customerEntity.ID, &customerEntity.Username, &customerEntity.Password,
		&customerEntity.Name, &customerEntity.Role, &customerEntity.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("find one customer by id: %w", common.ErrNotFound)
		}
		return nil, fmt.Errorf("find one customer by id: %w: %w", err, common.ErrInfrastructure)
	}

	return &customerEntity, nil
}

func (r *repo) FindOneByUsername(c context.Context, username string) (*customer.Customer, error) {
	var customerEntity customer.Customer

	sql := "SELECT id, username, password, name, role, created_at FROM customers.customers WHERE username = $1"
	row := r.db.QueryRow(c, sql, username)

	err := row.Scan(
		&customerEntity.ID, &customerEntity.Username, &customerEntity.Password,
		&customerEntity.Name, &customerEntity.Role, &customerEntity.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("find one customer by username %s: %w", username, common.ErrNotFound)
		}
		return nil, fmt.Errorf("find one customer by username %s: %w: %w", username, err, common.ErrInfrastructure)
	}

	return &customerEntity, nil
}

func (r *repo) Save(c context.Context, customer *customer.Customer) error {
	sql := `INSERT INTO customers.customers (id, username, password, name, role, created_at)
					VALUES ($1, $2, $3, $4, $5, $6)
					ON CONFLICT (id)
					DO UPDATE SET
						username = EXCLUDED.username,
						password = EXCLUDED.password,
						name = EXCLUDED.name,
						role = EXCLUDED.role`
	_, err := r.db.Exec(
		c, sql,
		customer.ID, customer.Username, customer.Password,
		customer.Name, customer.Role, customer.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("save customer: %w, %w", err, common.ErrInfrastructure)
	}

	return nil
}
