package session

import (
	"api/pkg/helpers"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	AccessToken string    `db:"access_token" json:"access_token"`
	CustomerID  uuid.UUID `db:"customer_id" json:"customer_id"`
	IP          string    `db:"ip" json:"-"`
	UserAgent   string    `db:"user_agent" json:"-"`
	CreatedAt   time.Time `db:"created_at" json:"-"`
}

func New(customerID uuid.UUID, ip, userAgent string) (*Session, error) {
	accessToken, err := helpers.GenerateRandomString(128)
	if err != nil {
		return nil, fmt.Errorf("session: failed to generate access token: %w", err)
	}

	return &Session{
		AccessToken: accessToken,
		CustomerID:  customerID,
		IP:          ip,
		UserAgent:   userAgent,
		CreatedAt:   time.Now().UTC(),
	}, nil
}
