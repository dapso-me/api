package customer

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Customer struct {
	ID        uuid.UUID `json:"id"`
	Login     string    `json:"login"`
	Password  string    `json:"-"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func New(login, password, name string) (*Customer, error) {
	ID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("new customer: failed to generate UUID: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, fmt.Errorf("new customer: failed to generate hash from pass: %w", err)
	}

	return &Customer{
		ID:        ID,
		Login:     login,
		Password:  string(hashedPassword),
		Name:      name,
		CreatedAt: time.Now().UTC(),
	}, nil
}
