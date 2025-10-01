package customer

import (
	"api/internal/common"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Role string

const (
	CustomerRole = "customer"
	OwnerRole    = "owner"
)

type Customer struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	Name      string    `json:"name"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

func New(username, password, name string) (*Customer, error) {
	if len(username) < 3 {
		return nil, fmt.Errorf("new customer: %w", ErrTooShortUsername)
	}

	if len(password) < 8 {
		return nil, fmt.Errorf("new customer: %w", ErrTooShortPassword)
	}

	if len(name) < 3 {
		return nil, fmt.Errorf("new customer: %w", ErrTooShortName)
	}

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
		Username:  username,
		Password:  string(hashedPassword),
		Name:      name,
		Role:      CustomerRole,
		CreatedAt: time.Now().UTC(),
	}, nil
}

func (c *Customer) ComparePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(c.Password), []byte(password))
	if err != nil {
		return fmt.Errorf("compare password: %w: %w", ErrIncorrectPassword, err)
	}

	return nil
}

func (c *Customer) HasRole(role Role) error {
	if c.Role != role {
		return fmt.Errorf("has role: %w", common.ErrForbidden)
	}
	return nil
}

func (c *Customer) UpdatePassword(currentPass, newPass string) error {
	if err := c.ComparePassword(currentPass); err != nil {
		return err
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPass), 10)
	if err != nil {
		return err
	}

	c.Password = string(hashedNewPassword)

	return nil
}

func (c *Customer) SetName(name string) {
	c.Name = name
}

func (c *Customer) SetUsername(username string) {
	c.Username = username
}
