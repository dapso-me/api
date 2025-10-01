package customer

import "api/internal/common"

var (
	ErrTooShortUsername  = common.NewAppError("TOO_SHORT_LOGIN", "the username is too short")
	ErrTooShortPassword  = common.NewAppError("TOO_SHORT_PASSWORD", "the password is too short")
	ErrTooShortName      = common.NewAppError("TOO_SHORT_NAME", "the name is too short")
	ErrIncorrectPassword = common.NewAppError("INCORRECT_PASSWORD", "the password is incorrect")
)
