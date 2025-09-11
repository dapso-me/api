package otp

import "api/internal/common"

var (
	ErrEmailRegistrationRateLimit = common.NewAppError("EMAIL_REGISTRATION_RATE_LIMIT", "registration limit exceeded for this email")
	ErrIncorrectOTP               = common.NewAppError("INCORRECT_OTP", "incorrect otp")
)
