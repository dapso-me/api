package session

import (
	"context"

	"github.com/google/uuid"
)

type UseCase interface {
	Create(c context.Context, dto *CreateInput) (*Session, error)
	RemoveOne(c context.Context, accessToken string) error
	RemoveAll(c context.Context, customerID uuid.UUID) error
}

type CreateInput struct {
	CustomerID uuid.UUID
	IP         string
	UserAgent  string
}
