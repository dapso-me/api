package menu_handler

import (
	"api/internal/domain/menu"

	"github.com/google/uuid"
)

type addCategoryReq struct {
	Position     int               `json:"position" example:"1"`
	Translations map[string]string `json:"translations" example:"{en:Category Name}"`
}

type addDishReq struct {
	CategoryID   uuid.UUID                       `json:"category_id"`
	Position     int                             `json:"position"`
	Photo        string                          `json:"photo"`
	Price        int                             `json:"price"`
	Translations map[string]menu.DishTranslation `json:"translations"`
}
