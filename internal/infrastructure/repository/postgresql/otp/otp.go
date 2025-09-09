package otp_postgresql_repository

import (
	"api/internal/common"
	"api/internal/domain/otp"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

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

func (r *repo) FindOneByCredentials(c context.Context, email, code string) (*otp.OTP, error) {
	query := "SELECT * FROM customers.otp WHERE email = $1 AND code = $2;"

	var otpEntity otp.OTP

	if err := r.db.GetContext(c, &otpEntity, query, email, code); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: find one by credentials: %w", common.ErrNotFound, err)
		}
		return nil, fmt.Errorf("%w: find one by credentials: %w", common.ErrInfrastructure, err)
	}

	return &otpEntity, nil
}

func (r *repo) FindByEmailSince(c context.Context, email string, since time.Time) ([]*otp.OTP, error) {
	query := "SELECT * FROM customers.otp WHERE email = $1 AND created_at > $2;"

	var otpEntities []*otp.OTP

	if err := r.db.SelectContext(c, &otpEntities, query, email, since); err != nil {
		return nil, fmt.Errorf("%w: find by email since: %w", common.ErrInfrastructure, err)
	}

	return otpEntities, nil
}

func (r *repo) FindByIPSince(c context.Context, IP string, since time.Time) ([]*otp.OTP, error) {
	query := "SELECT * FROM customers.otp WHERE ip = $1 AND created_at > $2;"

	var otpEntities []*otp.OTP

	if err := r.db.SelectContext(c, &otpEntities, query, IP, since); err != nil {
		return nil, fmt.Errorf("%w: find by IP since: %w", common.ErrInfrastructure, err)
	}

	return otpEntities, nil
}

func (r *repo) Save(c context.Context, otp *otp.OTP) error {
	query := `INSERT INTO customers.otp (id, email, code, purpose, ip, expires_at, created_at)
						VALUES ($1, $2, $3, $4, $5, $6, $7)
						ON CONFLICT (id) DO UPDATE
						SET email = $2, code = $3, purpose = $4, ip = $5, expires_at = $6, created_at = $7;`

	_, err := r.db.ExecContext(c, query, otp.ID, otp.Email, otp.Code, otp.Purpose, otp.IP, otp.ExpiresAt, otp.CreatedAt)
	if err != nil {
		return fmt.Errorf("%w: save: %w", common.ErrInfrastructure, err)
	}

	return nil
}

func (r *repo) Remove(c context.Context, otp *otp.OTP) error {
	query := "DELETE FROM customers.otp WHERE id = $1;"

	_, err := r.db.ExecContext(c, query, otp.ID)
	if err != nil {
		return fmt.Errorf("%w: remove: %w", common.ErrInfrastructure, err)
	}

	return nil
}
