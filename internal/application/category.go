package application

import (
	"api/internal/domain/menu"
	"context"

	"github.com/google/uuid"
)

type Category interface {
	FindByProjectID(c context.Context, projectID uuid.UUID) ([]*menu.Category, error)
	Add(c context.Context, input *AddCategoryInput) (*menu.Category, error)
	Update(c context.Context, input *UpdateCategoryInput) (*menu.Category, error)
	Remove(c context.Context, categoryID uuid.UUID) error
}

type AddCategoryInput struct {
	ProjectID    uuid.UUID
	Position     int
	Translations map[string]string
}

type UpdateCategoryInput struct {
}
