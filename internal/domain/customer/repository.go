package customer

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	FindOneByID(c context.Context, customerID uuid.UUID) (*Customer, error)
	FindOneByEmail(c context.Context, email string) (*Customer, error)

	Save(c context.Context, customer *Customer) error
	Remove(c context.Context, customer *Customer) error
}
