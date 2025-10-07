package project_usecase

import (
	"api/internal/application"
	"api/internal/common"
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
	auth        Auth
	projectRepo project.Repository
}

func New(auth Auth, projectRepo project.Repository) application.Project {
	return &uc{
		auth:        auth,
		projectRepo: projectRepo,
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
	if err := uc.auth.AdminGuard(c); err != nil {
		return nil, fmt.Errorf("add project: %w", err)
	}

	_, err := uc.projectRepo.FindOneBySlug(c, input.Slug)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			// it's okay
		} else {
			return nil, fmt.Errorf("add project: %w", err)
		}
	} else {
		return nil, fmt.Errorf("add project: %w", project.ErrSlugAlreadyInUse)
	}

	projectEntity, err := project.New(input.CustomerID, input.Slug, input.Name)
	if err != nil {
		return nil, fmt.Errorf("add project: %w", err)
	}

	if err := uc.projectRepo.Save(c, projectEntity); err != nil {
		return nil, fmt.Errorf("add project: %w", err)
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
