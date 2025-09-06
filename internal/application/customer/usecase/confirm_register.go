package customer_usecase

import (
	customer "api/internal/application/customer/domain"
	session "api/internal/application/session/domain"
	"api/internal/common"
	"context"
	"errors"
	"fmt"
)

func (uc *uc) ConfirmRegister(c context.Context, email, name, password, code string) (*session.Session, error) {
	{
		_, err := uc.customerRepo.FindOneByEmail(c, email)

		if err == nil {
			return nil, fmt.Errorf("register: %w", customer.ErrEmailAlreadyInUse)
		}

		if !errors.Is(err, common.ErrNotFound) {
			return nil, fmt.Errorf("register: %w", err)
		}
	}

	if err := uc.otpUC.Verify(c, email, code); err != nil {
		return nil, fmt.Errorf("register: %w", err)
	}

	cstmr, err := customer.New(email, password, name)
	if err != nil {
		return nil, fmt.Errorf("register: %w", err)
	}

	if err := uc.customerRepo.Save(c, cstmr); err != nil {
		return nil, fmt.Errorf("register: %w", err)
	}

	session, err := uc.sessionUC.Create(c, cstmr.ID)
	if err != nil {
		return nil, fmt.Errorf("register: %w", err)
	}

	return session, nil
}
