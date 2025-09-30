package customer

import "api/internal/common"

var (
	ErrTooShortLogin    = common.NewAppError("TOO_SHORT_LOGIN", "the login is too short")
	ErrTooShortPassword = common.NewAppError("TOO_SHORT_PASSWORD", "the password is too short")
	ErrTooShortName     = common.NewAppError("TOO_SHORT_NAME", "the name is too short")
)
