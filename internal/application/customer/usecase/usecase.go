package customer_usecase

import (
	customer "api/internal/application/customer/domain"
	session "api/internal/application/session/domain"
	"context"

	"github.com/google/uuid"
)

type OTPUseCase interface {
	Send(c context.Context, email, purpose string) error
	Verify(c context.Context, email, code string) error
}

type sessionUseCase interface {
	Create(c context.Context, customerID uuid.UUID) (*session.Session, error)
}

type uc struct {
	customerRepo customer.Repository
	otpUC        OTPUseCase
	sessionUC    sessionUseCase
}

func New(
	customerRepo customer.Repository,
	otpUC OTPUseCase,
	sessionUC sessionUseCase) *uc {
	return &uc{
		customerRepo: customerRepo,
		otpUC:        otpUC,
		sessionUC:    sessionUC,
	}
}
