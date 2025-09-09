package session_postgresql_repository

import (
	"api/internal/common"
	"api/internal/domain/session"
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

func (r *repo) FindOneByAccessToken(c context.Context, accessToken string) (*session.Session, error) {
	query := "SELECT * FROM customers.sessions WHERE access_token = $1;"

	var sessionEntity session.Session

	if err := r.db.GetContext(c, &sessionEntity, query, accessToken); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: find one by access token: %w", common.ErrNotFound, err)
		}

		return nil, fmt.Errorf("%w: find one by access token: %w", common.ErrInfrastructure, err)
	}

	return &sessionEntity, nil
}

func (r *repo) FindByCustomerID(c context.Context, customerID uuid.UUID) ([]*session.Session, error) {
	query := "SELECT * FROM customers.sessions WHERE customer_id = $1;"

	var sessionEntities []*session.Session

	if err := r.db.SelectContext(c, &sessionEntities, query, customerID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: find one by access token: %w", common.ErrNotFound, err)
		}

		return nil, fmt.Errorf("%w: find one by access token: %w", common.ErrInfrastructure, err)
	}

	return sessionEntities, nil
}

func (r *repo) Save(c context.Context, session *session.Session) error {
	query := `INSERT INTO customers.sessions (access_token, customer_id, ip, user_agent, created_at)
						VALUES ($1, $2, $3, $4, $5) ON CONFLICT (access_token)
						DO UPDATE ip = $3, user_agent = $4;`
	_, err := r.db.ExecContext(c, query, session.AccessToken, session.CustomerID, session.IP, session.UserAgent, session.CreatedAt)
	if err != nil {
		return fmt.Errorf("%w: save: %w", common.ErrInfrastructure, err)
	}

	return nil
}

func (r *repo) Remove(c context.Context, session *session.Session) error {
	query := "DELETE FROM customers.sessions WHERE access_token = $1;"

	_, err := r.db.ExecContext(c, query, session.AccessToken)
	if err != nil {
		return fmt.Errorf("%w: remove: %w", common.ErrInfrastructure, err)
	}

	return nil
}
