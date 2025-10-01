package session

import "context"

type Repository interface {
	FindOneByAccessToken(c context.Context, accessToken string) (*Session, error)
	Save(c context.Context, session *Session) error
}
