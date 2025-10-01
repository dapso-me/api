package project

import "api/internal/common"

var (
	ErrSlugAlreadyInUse = common.NewAppError("PROJECT_SLUG_ALREADY_IN_USE", "the slug already in use")
)
