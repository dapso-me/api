package application

import (
	"api/internal/domain/menu"
	"context"

	"github.com/google/uuid"
)

type Dish interface {
	Add(c context.Context, input *AddDishInput) (*menu.Dish, error)
	PutImage(c context.Context, input *PutDishImageInput) (*menu.Dish, error)
	Remove(c context.Context, dishID uuid.UUID) error
	Update(c context.Context, input *UpdateDishInput) (*menu.Dish, error)

	// AddOption(c context.Context, input *AddOptionInput) error
	// RemoveOption(c context.Context)
	// UpdateOption()

	// AddItem()
	// RemoveItem()
	// UpdateItem()
}

type AddDishInput struct {
	ProjectID    uuid.UUID
	CategoryID   uuid.UUID
	Position     int
	Photo        []byte
	Price        int
	Translations map[string]menu.DishTranslation
}

type PutDishImageInput struct {
	DishID uuid.UUID
	Photo  []byte
}

type UpdateDishInput struct {
	DishID       uuid.UUID
	CategoryID   *uuid.UUID
	Position     int
	Price        int
	Translations map[string]menu.DishTranslation
}
