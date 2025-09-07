package session

import (
	"api/internal/common"
	"context"
)

type contextKey string

const sessionKey contextKey = "session"

func GetFromContext(c context.Context) (*Session, error) {
	session, ok := c.Value(sessionKey).(*Session)
	if !ok {
		return nil, common.ErrUnauthorized
	}

	return session, nil
}

func SetToContext(c context.Context, s *Session) context.Context {
	return context.WithValue(c, sessionKey, s)
}
