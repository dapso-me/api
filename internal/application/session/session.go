package session_usecase

import (
	"api/internal/domain/session"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type useCase struct {
	sessionRepo session.Repository
}

func New(sessionRepo session.Repository) session.UseCase {
	return &useCase{
		sessionRepo: sessionRepo,
	}
}

func (uc *useCase) Create(c context.Context, dto *session.CreateInput) (*session.Session, error) {
	sessionEntity, err := session.New(dto.CustomerID, dto.IP, dto.UserAgent)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	if err := uc.sessionRepo.Save(c, sessionEntity); err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	return sessionEntity, nil
}

func (uc *useCase) RemoveOne(c context.Context, accessToken string) error {
	sessionEntity, err := uc.sessionRepo.FindOneByAccessToken(c, accessToken)
	if err != nil {
		return fmt.Errorf("remove one: %w", err)
	}

	if err := uc.sessionRepo.Remove(c, sessionEntity); err != nil {
		return fmt.Errorf("remove one: %w", err)
	}

	return nil
}

func (uc *useCase) RemoveAll(c context.Context, customerID uuid.UUID) error {
	if err := uc.RemoveAll(c, customerID); err != nil {
		return fmt.Errorf("remove all: %w", err)
	}

	return nil
}
