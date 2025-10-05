package menu

import (
	"context"

	"github.com/google/uuid"
)

type CategoryRepository interface {
	FindOneByID(c context.Context, categoryID uuid.UUID) (*Category, error)
	FindByProjectID(c context.Context, projectID uuid.UUID) ([]*Category, error)

	Save(c context.Context, category *Category) error
	Remove(c context.Context, category *Category) error
}

type DishRepository interface {
	FindOneByID(c context.Context, dishID uuid.UUID) (*Dish, error)
	FindByProjectID(c context.Context, projectID uuid.UUID) ([]*Dish, error)

	Save(c context.Context, dish *Dish) error
	Remove(c context.Context, dish *Dish) error
}
