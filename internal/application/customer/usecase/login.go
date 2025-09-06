package customer_usecase

import (
	session "api/internal/application/session/domain"
	"context"
	"fmt"
)

func (uc *uc) Login(c context.Context, email, password string) (*session.Session, error) {
	cstmr, err := uc.customerRepo.FindOneByEmail(c, email)
	if err != nil {
		return nil, fmt.Errorf("login: find customer by email: %w", err)
	}

	if err := cstmr.ComparePassword(password); err != nil {
		return nil, fmt.Errorf("login: invalid password: %w", err)
	}

	session, err := uc.sessionUC.Create(c, cstmr.ID)
	if err != nil {
		return nil, fmt.Errorf("login: failed to create session: %w", err)
	}

	return session, nil
}
