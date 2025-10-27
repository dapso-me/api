package project_usecase

import (
	"api/internal/application"
	"api/internal/common"
	"api/internal/domain/menu"
	"api/internal/domain/project"
	"api/internal/domain/session"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type Auth interface {
	AdminGuard(c context.Context) error
}

type uc struct {
	auth         Auth
	projectRepo  project.Repository
	categoryRepo menu.CategoryRepository
}

func New(auth Auth, projectRepo project.Repository, categoryRepo menu.CategoryRepository) application.Project {
	return &uc{
		auth:         auth,
		projectRepo:  projectRepo,
		categoryRepo: categoryRepo,
	}
}

func (uc *uc) ProjectGuard(c context.Context, projectID uuid.UUID) error {
	sessionEntity, err := session.GetFromContext(c)
	if err != nil {
		return fmt.Errorf("project guard: %w", err)
	}

	projectEntity, err := uc.projectRepo.FindOneByID(c, projectID)
	if err != nil {
		return fmt.Errorf("project guard: %w", err)
	}

	if projectEntity.CustomerID != sessionEntity.CustomerID {
		if err := uc.auth.AdminGuard(c); err != nil {
			return fmt.Errorf("project guard: %w", common.ErrForbidden)
		}
		return fmt.Errorf("project guard: %w", common.ErrForbidden)
	}

	return nil
}

func (uc *uc) FindOneBySlug(c context.Context, slug string) (*project.Project, error) {
	project, err := uc.projectRepo.FindOneBySlug(c, slug)
	if err != nil {
		return nil, fmt.Errorf("find one by slug: %w", err)
	}

	if project.IsRemoved() {
		return nil, fmt.Errorf("find one by slug: %w", common.ErrNotFound)
	}

	return project, nil
}

func (uc *uc) Add(c context.Context, input *application.AddProjectInput) (*project.Project, error) {
	wrapErr := func(err error) error {
		return fmt.Errorf("add project: %w", err)
	}

	if err := uc.auth.AdminGuard(c); err != nil {
		return nil, wrapErr(err)
	}

	_, err := uc.projectRepo.FindOneBySlug(c, input.Slug)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			// it's okay
		} else {
			return nil, wrapErr(err)
		}
	} else {
		return nil, wrapErr(project.ErrSlugAlreadyInUse)
	}

	projectEntity, err := project.New(input.CustomerID, input.Slug, input.Name)
	if err != nil {
		return nil, wrapErr(err)
	}

	if err := uc.projectRepo.Save(c, projectEntity); err != nil {
		return nil, wrapErr(err)
	}

	return projectEntity, nil
}

func (uc *uc) Remove(c context.Context, projectID uuid.UUID) error {
	if err := uc.auth.AdminGuard(c); err != nil {
		return fmt.Errorf("add project: %w", err)
	}

	projectEntity, err := uc.projectRepo.FindOneByID(c, projectID)
	if err != nil {
		return fmt.Errorf("add project: %w", err)
	}

	projectEntity.Remove()

	if err := uc.projectRepo.Save(c, projectEntity); err != nil {
		return fmt.Errorf("add project: %w", err)
	}

	return nil
}

func (uc *uc) FindByOwner(c context.Context) ([]*application.FullProject, error) {
	session, err := session.GetFromContext(c)
	if err != nil {
		return nil, fmt.Errorf("find by owner: %w", err)
	}

	output := []*application.FullProject{}

	projects, err := uc.projectRepo.FindByCustomerID(c, session.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("find by owner: %w", err)
	}

	for _, project := range projects {
		categories, err := uc.categoryRepo.FindByProjectID(c, project.ID)
		if err != nil {
			return nil, fmt.Errorf("find by owner: %w", err)
		}

		output = append(output, &application.FullProject{
			Project: project,
			Menu:    categories,
		})
	}

	return output, nil
}
