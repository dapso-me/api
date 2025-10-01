package application

import (
	"api/internal/domain/customer"
	"context"
)

type Customer interface {
	SignIn(c context.Context, input *SignInInput) (*AuthOutput, error)
	Authenticate(c context.Context) (*AuthOutput, error)
	ChangePassword(c context.Context, input *ChangePasswordInput) error
}

type AuthOutput struct {
	AccessToken string             `json:"access_token"`
	Customer    *customer.Customer `json:"customer"`
}

type SignInInput struct {
	Username  string
	Password  string
	IP        string
	UserAgent string
}

type ChangePasswordInput struct {
	CurrentPassword string
	NewPassword     string
}
