package otp

import "context"

type UseCase interface {
	Send(c context.Context, dto *SendInput) error
	VerifyAndRevoke(c context.Context, dto *VerifyInput) error
}

type SendInput struct {
	Email     string
	Purpose   Purpose
	IP        string
	UserAgent string
}

type VerifyInput struct {
	Email string
	Code  string
}
