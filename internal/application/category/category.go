package category_usecase

import (
	"api/internal/application"
	"api/internal/domain/menu"
	"api/internal/domain/project"
	"api/internal/domain/session"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type uc struct {
	categoryRepo menu.CategoryRepository
	dishRepo     menu.DishRepository
	projectRepo  project.Repository
}

func New(
	categoryRepo menu.CategoryRepository,
	projectRepo project.Repository,
	dishRepo menu.DishRepository) application.Category {
	return &uc{
		categoryRepo: categoryRepo,
		projectRepo:  projectRepo,
		dishRepo:     dishRepo,
	}
}

func (uc *uc) FindByProjectID(c context.Context, projectID uuid.UUID) ([]*menu.Category, error) {
	categories, err := uc.categoryRepo.FindByProjectID(c, projectID)
	if err != nil {
		return nil, fmt.Errorf("find by project ID: %w", err)
	}

	output, err := uc.loadDishes(c, projectID, categories)
	if err != nil {
		return nil, fmt.Errorf("find by project ID: %w", err)
	}

	return output, nil
}

func (uc *uc) loadDishes(c context.Context, projectID uuid.UUID, categories []*menu.Category) ([]*menu.Category, error) {
	dishes, err := uc.dishRepo.FindByProjectID(c, projectID)
	if err != nil {
		return nil, fmt.Errorf("loadDishes: %w", err)
	}

	for _, dish := range dishes {
		for _, category := range categories {
			category.Dishes = append(category.Dishes, *dish)
		}
	}

	return categories, nil
}

func (uc *uc) Add(c context.Context, input *application.AddCategoryInput) (*menu.Category, error) {
	session, err := session.GetFromContext(c)
	if err != nil {
		return nil, fmt.Errorf("add category: %w", err)
	}

	project, err := uc.projectRepo.FindOneByID(c, input.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("add category: %w", err)
	}

	if err := project.IsOwner(session.CustomerID); err != nil {
		return nil, fmt.Errorf("add category: %w", err)
	}

	categoryEntity, err := menu.NewCategory(input.ProjectID, input.Position, input.Translations)
	if err != nil {
		return nil, fmt.Errorf("add category: %w", err)
	}

	if err := uc.categoryRepo.Save(c, categoryEntity); err != nil {
		return nil, fmt.Errorf("add category: %w", err)
	}

	return categoryEntity, nil
}

func (uc *uc) Update(c context.Context, input *application.UpdateCategoryInput) (*menu.Category, error) {
	return nil, nil
}

func (uc *uc) Remove(c context.Context, categoryID uuid.UUID) error {
	session, err := session.GetFromContext(c)
	if err != nil {
		return fmt.Errorf("remove category: %w", err)
	}

	categoryEntity, err := uc.categoryRepo.FindOneByID(c, categoryID)
	if err != nil {
		return fmt.Errorf("remove category: %w", err)
	}

	projectEntity, err := uc.projectRepo.FindOneByID(c, categoryEntity.ProjectID)
	if err != nil {
		return fmt.Errorf("remove category: %w", err)
	}

	if err := projectEntity.IsOwner(session.CustomerID); err != nil {
		return fmt.Errorf("remove category: %w", err)
	}

	if err := uc.categoryRepo.Remove(c, categoryEntity); err != nil {
		return fmt.Errorf("remove catgory: %w", err)
	}

	return nil
}
