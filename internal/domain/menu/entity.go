package menu

import (
	"fmt"

	"github.com/google/uuid"
)

type Category struct {
	ID        uuid.UUID `json:"ID"`
	ProjectID uuid.UUID `json:"project_id"`
	Position  int       `json:"position"`
	Dishes    []Dish    `json:"dishes"`

	// key is lang_code
	Translations map[string]string `json:"translations"`
}

func NewCategory(projectID uuid.UUID, position int, translations map[string]string) (*Category, error) {
	ID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("new category: %w", err)
	}

	return &Category{
		ID:           ID,
		ProjectID:    projectID,
		Position:     position,
		Translations: translations,
	}, nil
}

type Dish struct {
	ID          uuid.UUID     `json:"ID"`
	ProjectID   uuid.UUID     `json:"project_id"`
	CategoryID  *uuid.UUID    `json:"category_id"`
	Position    int           `json:"position"`
	PhotoURL    string        `json:"photo_url"`
	Price       int           `json:"price"`
	IsAvailable bool          `json:"is_available"`
	Options     []OptionGroup `json:"options"`

	// key is lang_code
	Translations map[string]DishTranslation `json:"translations"`
}

type DishTranslation struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func NewDish(
	projectID uuid.UUID, position int, photoURL string,
	price int, translations map[string]DishTranslation) (*Dish, error) {

	ID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("new dish: %w", err)
	}

	return &Dish{
		ID:           ID,
		ProjectID:    projectID,
		CategoryID:   nil,
		Position:     position,
		PhotoURL:     photoURL,
		Price:        price,
		IsAvailable:  true,
		Translations: translations,
	}, nil
}

func (d *Dish) SetPhoto(photoURL string) {
	d.PhotoURL = photoURL
}

func (d *Dish) Update(categoryID *uuid.UUID, position int, price int, translations map[string]DishTranslation) {
	d.CategoryID = categoryID
	d.Position = position
	d.Price = price
	d.Translations = translations
}

type OptionGroup struct {
	ID           uuid.UUID         `json:"ID"`
	DishID       uuid.UUID         `json:"dish_id"`
	Position     int               `json:"position"`
	MinQuantity  int               `json:"min_quantity"`
	MaxQuantity  int               `json:"max_quantity"`
	Translations map[string]string `json:"translations"`
	Items        []OptionItem      `json:"items"`
}

func NewOptionGroup(dishID uuid.UUID, position, minQuantity, maxQuantity int, translations map[string]string) (*OptionGroup, error) {
	ID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("new option group: %w", err)
	}

	if minQuantity > maxQuantity {
		return nil, fmt.Errorf("new option group: %w", ErrValueError)
	}

	return &OptionGroup{
		ID:           ID,
		DishID:       dishID,
		Position:     position,
		MinQuantity:  minQuantity,
		MaxQuantity:  maxQuantity,
		Translations: translations,
	}, nil
}

type OptionItem struct {
	ID            uuid.UUID         `json:"ID"`
	OptionGroupID uuid.UUID         `json:"option_group_id"`
	Position      int               `json:"position"`
	MinQuantity   int               `json:"min_quantity"`
	MaxQuantity   int               `json:"max_quantity"`
	PriceModifier int               `json:"price_modifier"`
	Translations  map[string]string `json:"translations"`
}

func NewOptionItem(
	optionGroupID uuid.UUID,
	position, minQuantity, maxQuantity, priceModifier int,
	translations map[string]string,
) (*OptionItem, error) {

	ID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("new option item: %w", err)
	}

	if minQuantity > maxQuantity {
		return nil, fmt.Errorf("new option group: %w", ErrValueError)
	}

	return &OptionItem{
		ID:            ID,
		OptionGroupID: optionGroupID,
		Position:      position,
		MinQuantity:   minQuantity,
		MaxQuantity:   maxQuantity,
		PriceModifier: priceModifier,
		Translations:  translations,
	}, nil
}
