package otp

import (
	"api/pkg/helpers"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Purpose string

const (
	RegisterPurpose Purpose = "register"
	RecoveryPurpose Purpose = "recovery"
)

const TTL = time.Second * 60 * 10

type OTP struct {
	ID        uuid.UUID `db:"id"`
	Email     string    `db:"email"`
	Code      string    `db:"code"`
	Purpose   Purpose   `db:"purpose"`
	IP        string    `db:"ip"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}

func New(email string, purpose Purpose, IP string) (*OTP, error) {
	ID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("otp: failed to generate UUID: %w", err)
	}

	code, err := helpers.GenerateOTP(6)
	if err != nil {
		return nil, fmt.Errorf("otp: failed to generate code")
	}

	currentTime := time.Now().UTC()

	return &OTP{
		ID:        ID,
		Email:     email,
		Code:      code,
		Purpose:   purpose,
		IP:        IP,
		ExpiresAt: currentTime.Add(TTL),
		CreatedAt: currentTime,
	}, nil
}
