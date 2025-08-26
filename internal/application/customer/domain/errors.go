package customer

import (
	"api/internal/common"
)

var (
	ErrInvalidEmail      = common.NewAppError("INVALID_EMAIL", "invalid email address")
	ErrNameEmpty         = common.NewAppError("EMPTY_NAME", "name can't be empty")
	ErrTooShortPassword  = common.NewAppError("TOO_SHORT_PASSWORD", "password must be 8 characters or more")
	ErrIncorrectPassword = common.NewAppError("INCORRECT_PASSWORD", "incorrect password")
)
