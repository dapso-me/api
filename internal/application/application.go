package application

import "context"

type Authorization interface {
	SignIn(c context.Context, input *SignInInput) (*AuthOutput, error)
	Authenticate(c context.Context) (*AuthOutput, error)
}

type AuthOutput struct {
	AccessToken string `json:"access_token"`
}

type SignInInput struct {
	Login     string
	Password  string
	Name      string
	IP        string
	UserAgent string
}

type Project interface {
	FindProjectBySlug(c context.Context, slug string) (*ProjectOutput, error)
}

type ProjectOutput struct {
}
