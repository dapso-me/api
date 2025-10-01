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
	IP          string    `json:"ip"`
	UserAgent   string    `json:"user_agent"`
	CreatedAt   time.Time `json:"created_at"`
}

func New(customerID uuid.UUID, ip, userAgent string) (*Session, error) {
	accessToken, err := helpers.GenerateRandomString(128)
	if err != nil {
		return nil, fmt.Errorf("new session: failed to generate random string: %w", err)
	}

	return &Session{
		AccessToken: accessToken,
		CustomerID:  customerID,
		IP:          ip,
		UserAgent:   userAgent,
		CreatedAt:   time.Now().UTC(),
	}, nil
}
