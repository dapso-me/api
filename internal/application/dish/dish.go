package dish_usecase

import (
	"api/internal/application"
	"api/internal/domain/menu"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type auth interface {
	ProjectGuard(c context.Context, projectID uuid.UUID) error
}

type storage interface {
	PutImage(c context.Context, object []byte, resizeWidth int) (string, error)
	Remove(c context.Context, objectName string) error
}

type uc struct {
	dishRepo     menu.DishRepository
	categoryRepo menu.CategoryRepository
	auth         auth
	storage      storage
}

func New(dishRepo menu.DishRepository, auth auth, storage storage) application.Dish {
	return &uc{
		dishRepo: dishRepo,
		auth:     auth,
		storage:  storage,
	}
}

func (uc *uc) Add(c context.Context, i *application.AddDishInput) (*menu.Dish, error) {
	if err := uc.auth.ProjectGuard(c, i.ProjectID); err != nil {
		return nil, fmt.Errorf("add dish: %w", err)
	}

	photoURL, err := uc.storage.PutImage(c, i.Photo, 500)
	if err != nil {
		return nil, fmt.Errorf("add dish: %w", err)
	}

	dishEntity, err := menu.NewDish(i.ProjectID, i.Position, photoURL, i.Price, i.Translations)
	if err != nil {
		return nil, fmt.Errorf("add dish: %w", err)
	}

	if err := uc.dishRepo.Save(c, dishEntity); err != nil {
		return nil, fmt.Errorf("add dish: %w", err)
	}

	return dishEntity, nil
}

func (uc *uc) PutImage(c context.Context, i *application.PutDishImageInput) (*menu.Dish, error) {
	dishEntity, err := uc.dishRepo.FindOneByID(c, i.DishID)
	if err != nil {
		return nil, fmt.Errorf("put image to dish %v: %w", i.DishID, err)
	}

	if err := uc.auth.ProjectGuard(c, dishEntity.ProjectID); err != nil {
		return nil, fmt.Errorf("put image to dish %v: %w", i.DishID, err)
	}

	oldPhoto := dishEntity.PhotoURL

	newPhotoURL, err := uc.storage.PutImage(c, i.Photo, 500)
	if err != nil {
		return nil, fmt.Errorf("put image to dish %v: %w", i.DishID, err)
	}

	dishEntity.SetPhoto(newPhotoURL)

	if err := uc.dishRepo.Save(c, dishEntity); err != nil {
		return nil, fmt.Errorf("put image to dish %v: %w", i.DishID, err)
	}

	if oldPhoto != "" {
		uc.storage.Remove(c, dishEntity.PhotoURL)
	}

	return dishEntity, nil
}

func (uc *uc) Remove(c context.Context, dishID uuid.UUID) error {
	dishEntity, err := uc.dishRepo.FindOneByID(c, dishID)
	if err != nil {
		return fmt.Errorf("remove dish %v: %w", dishID, err)
	}

	if err := uc.auth.ProjectGuard(c, dishEntity.ProjectID); err != nil {
		return fmt.Errorf("remove dish %v: %w", dishID, err)
	}

	if err := uc.dishRepo.Remove(c, dishEntity); err != nil {
		return fmt.Errorf("remove dish %v: %w", dishID, err)
	}

	return nil
}

func (uc *uc) Update(c context.Context, i *application.UpdateDishInput) (*menu.Dish, error) {
	dishEntity, err := uc.dishRepo.FindOneByID(c, i.DishID)
	if err != nil {
		return nil, fmt.Errorf("update dish %v: %w", i.DishID, err)
	}

	if err := uc.auth.ProjectGuard(c, dishEntity.ProjectID); err != nil {
		return nil, fmt.Errorf("update dish %v: %w", i.DishID, err)
	}

	dishEntity.Update(i.CategoryID, i.Position, i.Price, i.Translations)

	if err := uc.dishRepo.Save(c, dishEntity); err != nil {
		return nil, fmt.Errorf("update dish %v: %w", i.DishID, err)
	}

	return dishEntity, nil
}
