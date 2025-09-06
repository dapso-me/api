package session_usecase

import session "api/internal/application/session/domain"

type uc struct {
	sessionRepo session.Repository
}

func New(sessionRepo session.Repository) *uc {
	return &uc{
		sessionRepo: sessionRepo,
	}
}
