package application

import (
	"api/internal/domain/menu"
	"context"

	"github.com/google/uuid"
)

type Dish interface {
	Add(c context.Context, input *AddCategoryInput) (*menu.Dish, error)
	Remove(c context.Context, dishID uuid.UUID) error
	Update()

	AddOption(c context.Context, input *AddOptionInput) error
	RemoveOption(c context.Context)
	UpdateOption()

	AddItem()
	RemoveItem()
	UpdateItem()
}

type AddDishInput struct {
	CategoryID   uuid.UUID
	Position     int
	PhotoURL     string
	Price        int
	Translations map[string]menu.DishTranslation
}

type AddOptionInput struct {
	DishID uuid.UUID
}
