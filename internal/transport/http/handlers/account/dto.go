package account_handler

import "api/internal/application"

type signInReq struct {
	Username string `json:"username" example:"john_doe"`
	Password string `json:"password" example:"qwerty123"`
}

type changePasswordReq struct {
	CurrentPassword string `json:"current_password" example:"qwerty123"`
	NewPassword     string `json:"new_password" example:"321ytrewq"`
}

type AuthRes struct {
	*application.AuthOutput
	Projects []*application.FullProject `json:"projects"`
}
