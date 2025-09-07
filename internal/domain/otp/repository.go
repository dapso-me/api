package otp

import (
	"context"
	"time"
)

type Repository interface {
	FindOneByCredentials(c context.Context, email, code string) (*OTP, error)
	FindByEmailSince(c context.Context, email string, since time.Time) ([]*OTP, error)
	FindByIPSince(c context.Context, IP string, since time.Time) ([]*OTP, error)

	Save(c context.Context, otp *OTP) error
	Remove(c context.Context, otp *OTP) error
}
