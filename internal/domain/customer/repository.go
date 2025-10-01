package customer

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	FindOneByID(c context.Context, customerID uuid.UUID) (*Customer, error)
	FindOneByUsername(c context.Context, username string) (*Customer, error)

	Save(c context.Context, customer *Customer) error
}
