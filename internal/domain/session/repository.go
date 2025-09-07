package session

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	FindOneByAccessToken(c context.Context, accessToken string) (*Session, error)
	FindByCustomerID(c context.Context, customerID uuid.UUID) ([]*Session, error)
	Save(c context.Context, session *Session) error
	Remove(c context.Context, session *Session) error
}
