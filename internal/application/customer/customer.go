package customer_usecase

import (
	"api/internal/application"
	"api/internal/common"
	"api/internal/domain/customer"
	"api/internal/domain/session"
	"context"
	"errors"
	"fmt"
)

type uc struct {
	customerRepo customer.Repository
	sessionRepo  session.Repository
}

func New(customerRepo customer.Repository, sessionRepo session.Repository) application.Customer {
	return &uc{
		customerRepo: customerRepo,
		sessionRepo:  sessionRepo,
	}
}

func (uc *uc) AdminGuard(c context.Context) error {
	session, err := session.GetFromContext(c)
	if err != nil {
		return fmt.Errorf("admin guard: %w", err)
	}

	customerEntity, err := uc.customerRepo.FindOneByID(c, session.CustomerID)
	if err != nil {
		return fmt.Errorf("admin guard: %w", err)
	}

	if err := customerEntity.HasRole(customer.OwnerRole); err != nil {
		return fmt.Errorf("admin guard: %w", err)
	}

	return nil
}

func (uc *uc) SignIn(c context.Context, input *application.SignInInput) (*application.AuthOutput, error) {
	customerEntity, err := uc.customerRepo.FindOneByUsername(c, input.Username)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return nil, fmt.Errorf("sign in: %w", customer.ErrIncorrectPassword)
		}
		return nil, fmt.Errorf("sign in: %w", err)
	}

	if err := customerEntity.ComparePassword(input.Password); err != nil {
		return nil, fmt.Errorf("sign in: %w", err)
	}

	sessionEntity, err := session.New(customerEntity.ID, input.IP, input.UserAgent)
	if err != nil {
		return nil, fmt.Errorf("sign in: %w", err)
	}

	if err := uc.sessionRepo.Save(c, sessionEntity); err != nil {
		return nil, fmt.Errorf("sign in: %w", err)
	}

	return &application.AuthOutput{
		AccessToken: sessionEntity.AccessToken,
		Customer:    customerEntity,
	}, nil
}

func (uc *uc) Authenticate(c context.Context) (*application.AuthOutput, error) {
	sessionEntity, err := session.GetFromContext(c)
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	customerEntity, err := uc.customerRepo.FindOneByID(c, sessionEntity.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	return &application.AuthOutput{
		AccessToken: sessionEntity.AccessToken,
		Customer:    customerEntity,
	}, nil
}

func (uc *uc) ChangePassword(c context.Context, input *application.ChangePasswordInput) error {
	session, err := session.GetFromContext(c)
	if err != nil {
		return fmt.Errorf("change password: %w", err)
	}

	customerEntity, err := uc.customerRepo.FindOneByID(c, session.CustomerID)
	if err != nil {
		return fmt.Errorf("change password: %w", err)
	}

	if err := customerEntity.UpdatePassword(input.CurrentPassword, input.NewPassword); err != nil {
		return fmt.Errorf("change password: %w", err)
	}

	return nil
}
