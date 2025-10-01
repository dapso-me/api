package project_postgresql_repository

import (
	"api/internal/common"
	"api/internal/domain/project"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repo struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) project.Repository {
	return &repo{
		db: db,
	}
}

func (r *repo) FindOneBySlug(c context.Context, slug string) (*project.Project, error) {
	var projectEntity project.Project

	sql := "SELECT id, customer_id, slug, name, removed_at, created_at FROM project.projects WHERE slug = $1"
	row := r.db.QueryRow(c, sql, slug)

	err := row.Scan(
		&projectEntity.ID, &projectEntity.CustomerID, &projectEntity.Slug,
		&projectEntity.Name, &projectEntity.RemovedAt, &projectEntity.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("find one project by slug %s: %w: %w", slug, err, common.ErrNotFound)
		}
		return nil, fmt.Errorf("find one project by slug %s: %w: %w", slug, err, common.ErrInfrastructure)
	}

	return &projectEntity, nil
}

func (r *repo) FindOneByID(c context.Context, projectID uuid.UUID) (*project.Project, error) {
	var projectEntity project.Project

	sql := "SELECT id, customer_id, slug, name, removed_at, created_at FROM project.projects WHERE id = $1"
	row := r.db.QueryRow(c, sql, projectID)

	err := row.Scan(
		&projectEntity.ID, &projectEntity.CustomerID, &projectEntity.Slug,
		&projectEntity.Name, &projectEntity.RemovedAt, &projectEntity.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("find one project by ID %s: %w: %w", projectID, err, common.ErrNotFound)
		}
		return nil, fmt.Errorf("find one project by ID %s: %w: %w", projectID, err, common.ErrInfrastructure)
	}

	return &projectEntity, nil
}

func (r *repo) Save(c context.Context, project *project.Project) error {
	sql := `INSERT INTO project.projects (id, customer_id, slug, name, removed_at, created_at)
						VALUES ($1, $2, $3, $4, $5, $6)
					ON CONFLICT (id)
					DO UPDATE SET
						slug = EXCLUDED.slug,
						name = EXCLUDED.name,
						removed_at = EXCLUDED.removed_at`
	_, err := r.db.Exec(c, sql, project.ID, project.CustomerID, project.Slug, project.Name, project.RemovedAt, project.CreatedAt)
	if err != nil {
		return fmt.Errorf("save project: %w", err)
	}

	return nil
}
