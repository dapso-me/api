package session_postgresql_repository

import (
	"api/internal/common"
	"api/internal/domain/session"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repo struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) session.Repository {
	return &repo{
		db: db,
	}
}

func (r *repo) FindOneByAccessToken(c context.Context, accessToken string) (*session.Session, error) {
	var sessionEntity session.Session

	sql := "SELECT access_token, customer_id, ip, user_agent, created_at FROM customers.sessions WHERE access_token = $1"
	row := r.db.QueryRow(c, sql, accessToken)

	err := row.Scan(
		&sessionEntity.AccessToken, &sessionEntity.CustomerID,
		&sessionEntity.IP, &sessionEntity.UserAgent, &sessionEntity.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("find one session by access token: %w: %w", err, common.ErrNotFound)
		}
		return nil, fmt.Errorf("find one session by access token: %w: %w", err, common.ErrInfrastructure)
	}

	return &sessionEntity, nil
}

func (r *repo) Save(c context.Context, session *session.Session) error {
	sql := `INSERT INTO customers.sessions (access_token, customer_id, ip, user_agent, created_at)
					VALUES ($1, $2, $3, $4, $5)
					ON CONFLICT (access_token)
					DO UPDATE SET
						ip = EXCLUDED.ip,
						user_agent = EXCLUDED.user_agent`
	_, err := r.db.Exec(c, sql, session.AccessToken, session.CustomerID, session.IP, session.UserAgent, session.CreatedAt)
	if err != nil {
		return fmt.Errorf("save session: %w: %w", err, common.ErrInfrastructure)
	}

	return nil
}
