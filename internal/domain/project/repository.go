package project

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	FindOneBySlug(c context.Context, slug string) (*Project, error)
	FindOneByID(c context.Context, projectID uuid.UUID) (*Project, error)

	Save(c context.Context, project *Project) error
}
