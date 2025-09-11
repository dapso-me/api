package customer

import (
	"fmt"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Customer struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"-"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"-"`
}

func New(email, password, name string) (*Customer, error) {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return nil, fmt.Errorf("customer: %w", ErrInvalidEmail)
	}

	if len(password) < 8 {
		return nil, fmt.Errorf("customer: %w", ErrTooShortPassword)
	}

	if len(name) < 1 {
		return nil, fmt.Errorf("customer: %w", ErrNameEmpty)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, fmt.Errorf("customer: failed to generate hash from password: %w", err)
	}

	return &Customer{
		ID:        uuid.New(),
		Email:     email,
		Password:  string(hashedPassword),
		Name:      name,
		CreatedAt: time.Now().UTC(),
	}, nil
}

func (c *Customer) ComparePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(c.Password), []byte(password))
	if err != nil {
		return ErrIncorrectPassword
	}

	return nil
}

func (c *Customer) SetName(name string) {
	c.Name = name
}

func (c *Customer) SetPassword(currentPassword, newPassword string) error {
	if err := c.ComparePassword(currentPassword); err != nil {
		return err
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
	if err != nil {
		return fmt.Errorf("customer: failed to hash password: %w", err)
	}

	c.Password = string(hashedNewPassword)

	return nil
}
