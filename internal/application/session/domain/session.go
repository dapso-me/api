package session

import (
	"api/pkg/helpers"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	AccessToken string    `json:"access_token"`
	CustomerID  uuid.UUID `json:"customer_id"`
	IP          string    `json:"-"`
	UserAgent   string    `json:"-"`
	CreatedAt   time.Time `json:"-"`
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
