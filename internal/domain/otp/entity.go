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
	ID        uuid.UUID
	Email     string
	Code      string
	Purpose   Purpose
	IP        string
	ExpiresAt time.Time
	CreatedAt time.Time
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
