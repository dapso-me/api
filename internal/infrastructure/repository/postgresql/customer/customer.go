package customer_repository

import (
	"api/internal/common"
	"api/internal/domain/customer"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *repo {
	return &repo{
		db: db,
	}
}

func (r *repo) FindOneByID(c context.Context, ID uuid.UUID) (*customer.Customer, error) {
	var customerEntity customer.Customer

	query := `SELECT * FROM customers.customers WHERE id = $1;`

	err := r.db.GetContext(c, &customerEntity, query, ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: find one by id: %w", common.ErrNotFound, err)
		}

		return nil, fmt.Errorf("%w: find one by id: %w", common.ErrInfrastructure, err)
	}

	return &customerEntity, nil
}

func (r *repo) FindOneByEmail(c context.Context, email string) (*customer.Customer, error) {
	var customerEntity customer.Customer

	query := `SELECT * FROM customers.customers WHERE email = $1;`

	err := r.db.GetContext(c, &customerEntity, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: find one by id: %w", common.ErrNotFound, err)
		}

		return nil, fmt.Errorf("%w: find one by email: %w", common.ErrInfrastructure, err)
	}

	return &customerEntity, nil
}

func (r *repo) Save(c context.Context, customer *customer.Customer) error {
	query := `INSERT INTO customers.customers (id, email, password, name, created_at)
						VALUES ($1, $2, $3, $4, $5)
						ON CONFLICT (id) DO UPDATE SET password = $3, name = $4;`
	_, err := r.db.ExecContext(c, query, customer.ID, customer.Email, customer.Password, customer.Name, customer.CreatedAt)
	if err != nil {
		return fmt.Errorf("%w: save: %w", common.ErrInfrastructure, err)
	}

	return nil
}

func (r *repo) Remove(c context.Context, customer *customer.Customer) error {
	query := "DELETE FROM customers.customers WHERE id = $1;"

	_, err := r.db.ExecContext(c, query, customer.ID)
	if err != nil {
		return fmt.Errorf("%w: remove: %w", common.ErrInfrastructure, err)
	}

	return nil
}
