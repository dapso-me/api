package customer_usecase

import (
	customer "api/internal/application/customer/domain"
	"api/internal/common"
	"context"
	"errors"
	"fmt"
)

func (uc *uc) Register(c context.Context, email string) error {
	_, err := uc.customerRepo.FindOneByEmail(c, email)

	if err == nil {
		return fmt.Errorf("register: %w", customer.ErrEmailAlreadyInUse)
	}

	if !errors.Is(err, common.ErrNotFound) {
		return fmt.Errorf("register: %w", err)
	}

	if err := uc.otpUC.Send(c, email, "register"); err != nil {
		return fmt.Errorf("register: failed to send code: %w", err)
	}

	return nil
}
