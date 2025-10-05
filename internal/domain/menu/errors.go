package menu

import "api/internal/common"

var (
	ErrValueError = common.NewAppError("VALUE_ERROR", "minimum value cannot be greater than the maximum value")
)
