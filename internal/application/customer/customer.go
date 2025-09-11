package customer_usecase

import (
	"api/internal/common"
	"api/internal/domain/customer"
	"api/internal/domain/otp"
	"api/internal/domain/session"
	"context"
	"errors"
	"fmt"
)

type useCase struct {
	customerRepo customer.Repository

	otpService     otp.UseCase
	sessionService session.UseCase
}

func New(customerRepo customer.Repository, otpService otp.UseCase, sessionService session.UseCase) customer.UseCase {
	return &useCase{
		customerRepo:   customerRepo,
		otpService:     otpService,
		sessionService: sessionService,
	}
}

func (uc *useCase) IsEmailAvailable(c context.Context, email string) error {
	_, err := uc.customerRepo.FindOneByEmail(c, email)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return nil
		} else {
			return fmt.Errorf("emailAvailable: %w", err)
		}
	} else {
		return fmt.Errorf("emailAvailable: %w", customer.ErrEmailAlreadyInUse)
	}
}

func (uc *useCase) Register(c context.Context, dto *customer.RegisterInput) error {
	if err := uc.IsEmailAvailable(c, dto.Email); err != nil {
		return fmt.Errorf("Register: %w", err)
	}

	otpInput := &otp.SendInput{
		Email:     dto.Email,
		Purpose:   otp.RegisterPurpose,
		IP:        dto.IP,
		UserAgent: dto.UserAgent,
	}
	if err := uc.otpService.Send(c, otpInput); err != nil {
		return fmt.Errorf("Register: %w", err)
	}

	return nil
}

func (uc *useCase) ConfirmRegister(c context.Context, dto *customer.ConfirmRegisterInput) (*customer.AuthOutput, error) {
	otpInput := &otp.VerifyInput{
		Email: dto.Email,
		Code:  dto.Code,
	}
	if err := uc.otpService.VerifyAndRevoke(c, otpInput); err != nil {
		return nil, fmt.Errorf("Confirm Register: %w", err)
	}

	customerEntity, err := customer.New(dto.Email, dto.Password, dto.Name)
	if err != nil {
		return nil, fmt.Errorf("Confirm Register: %w", err)
	}

	if err := uc.customerRepo.Save(c, customerEntity); err != nil {
		return nil, fmt.Errorf("Confirm Register: %w", err)
	}

	sessionInput := &session.CreateInput{
		CustomerID: customerEntity.ID,
		IP:         dto.IP,
		UserAgent:  dto.UserAgent,
	}
	sessionEntity, err := uc.sessionService.Create(c, sessionInput)
	if err != nil {
		return nil, fmt.Errorf("Confirm Register: %w", err)
	}

	return &customer.AuthOutput{
		AccessToken: sessionEntity.AccessToken,
		Customer:    customerEntity,
	}, nil
}

func (uc *useCase) Login(c context.Context, dto *customer.LoginInput) (*customer.AuthOutput, error) {
	customerEntity, err := uc.customerRepo.FindOneByEmail(c, dto.Email)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return nil, fmt.Errorf("login: %w", customer.ErrIncorrectPassword)
		}

		return nil, fmt.Errorf("login: %w", err)
	}

	if err := customerEntity.ComparePassword(dto.Password); err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}

	createSessionInput := &session.CreateInput{
		CustomerID: customerEntity.ID,
		IP:         dto.IP,
		UserAgent:  dto.UserAgent,
	}
	session, err := uc.sessionService.Create(c, createSessionInput)
	if err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}

	return &customer.AuthOutput{
		AccessToken: session.AccessToken,
		Customer:    customerEntity,
	}, nil
}

func (uc *useCase) Authenticate(c context.Context) (*customer.AuthOutput, error) {
	session, err := session.GetFromContext(c)
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	customerEntity, err := uc.customerRepo.FindOneByID(c, session.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	return &customer.AuthOutput{
		AccessToken: session.AccessToken,
		Customer:    customerEntity,
	}, nil
}

func (uc *useCase) Update(c context.Context, dto *customer.UpdateInput) error {
	session, err := session.GetFromContext(c)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}

	customerEntity, err := uc.customerRepo.FindOneByID(c, session.CustomerID)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}

	customerEntity.SetName(dto.Name)

	if err := uc.customerRepo.Save(c, customerEntity); err != nil {
		return fmt.Errorf("update: %w", err)
	}

	return nil
}

func (uc *useCase) UpdatePassword(c context.Context, dto *customer.UpdatePasswordInput) error {
	session, err := session.GetFromContext(c)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	customerEntity, err := uc.customerRepo.FindOneByID(c, session.CustomerID)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	if err := customerEntity.SetPassword(dto.CurrentPassword, dto.NewPassword); err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	if err := uc.customerRepo.Save(c, customerEntity); err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	return nil
}
