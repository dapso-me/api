package session_usecase

import (
	session "api/internal/application/session/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (uc *uc) Create(c context.Context, customerID uuid.UUID, ip, userAgent string) (*session.Session, error) {
	session, err := session.New(customerID, ip, userAgent)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	if err := uc.sessionRepo.Save(c, session); err != nil {
		return nil, fmt.Errorf("create session: failed to save: %w", err)
	}

	return session, nil
}
