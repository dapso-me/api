package customer

import (
	"context"
)

type UseCase interface {
	Authenticate(c context.Context) (*AuthOutput, error)
	Register(c context.Context, dto *RegisterInput) error
	ConfirmRegister(c context.Context, dto *ConfirmRegisterInput) (*AuthOutput, error)
	Login(c context.Context, dto *LoginInput) (*AuthOutput, error)
	Update(c context.Context, dto *UpdateInput) error
	UpdatePassword(c context.Context, dto *UpdatePasswordInput) error
}

type AuthOutput struct {
	AccessToken string    `json:"access_token"`
	Customer    *Customer `json:"customer"`
}

type LoginInput struct {
	Email     string
	Password  string
	IP        string
	UserAgent string
}

type RegisterInput struct {
	Email     string
	IP        string
	UserAgent string
}

type ConfirmRegisterInput struct {
	Email     string
	Password  string
	Code      string
	Name      string
	IP        string
	UserAgent string
}

type UpdateInput struct {
	Name string
}

type UpdatePasswordInput struct {
	CurrentPassword string
	NewPassword     string
}
