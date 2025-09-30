package dish

import "github.com/google/uuid"

type Dish struct {
	ID           uuid.UUID `json:"id"`
	CategoryID   uuid.UUID `json:"category"`
	PhotoURL     string    `json:"photo_url"`
	PhotoMiniURL string    `json:"photo_mini_url"`
	Price        int       `json:"price"`
	IsAvailable  bool      `json:"is_available"`
	Options      []Option  `json:"options"`

	// map key is lang code like en, kz...
	Translations map[string]*DishTranslation `json:"translations"`
}

type DishTranslation struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Option struct {
	ID          uuid.UUID `json:"id"`
	MinQuantity int       `json:"min_quantity"`
	MaxQuantity int       `json:"max_quantity"`

	// map key is lang code like en, kz...
	Translations map[string]*OptionTranslation `json:"translations"`
}

type OptionTranslation struct {
	Name string `json:"name"`
}

type Item struct {
	ID            uuid.UUID `json:"id"`
	MinQuanity    int       `json:"min_quantity"`
	MaxQuantity   int       `json:"max_quantity"`
	PriceModifier int       `json:"price_modifier"`

	// map key is lang code like en, kz...
	Translations map[string]*ItemTranslation `json:"translations"`
}

type ItemTranslation struct {
	Name string `json:"name"`
}

func NewDish()
