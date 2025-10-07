package application

import (
	"api/internal/domain/project"
	"context"

	"github.com/google/uuid"
)

type Project interface {
	FindOneBySlug(c context.Context, slug string) (*project.Project, error)

	Add(c context.Context, input *AddProjectInput) (*project.Project, error)
	Remove(c context.Context, projectID uuid.UUID) error

	ProjectGuard(c context.Context, projectID uuid.UUID) error
}

type AddProjectInput struct {
	CustomerID uuid.UUID
	Slug       string
	Name       string
}
