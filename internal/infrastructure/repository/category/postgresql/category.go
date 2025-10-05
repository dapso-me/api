package category_postgresql_repository

import (
	"api/internal/common"
	"api/internal/domain/menu"
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

func New(db *pgxpool.Pool) menu.CategoryRepository {
	return &repo{
		db: db,
	}
}

func (r *repo) FindOneByID(c context.Context, categoryID uuid.UUID) (*menu.Category, error) {
	var categoryEntity menu.Category

	sql := `SELECT id, project_id, position FROM module_menu.categories WHERE id = $1`
	row := r.db.QueryRow(c, sql, categoryID)

	if err := row.Scan(&categoryEntity.ID, &categoryEntity.ProjectID, &categoryEntity.Position); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("find one category by ID %v: %w: %w", categoryID, err, common.ErrNotFound)
		}
		return nil, fmt.Errorf("find one category by ID %v: %w", categoryID, err)
	}

	sql = "SELECT lang_code, name FROM module_menu.category_translations WHERE category_id = $1;"
	rows, err := r.db.Query(c, sql, categoryID)
	if err != nil {
		return nil, fmt.Errorf("find one category by ID %v: %w", categoryID, err)
	}
	defer rows.Close()

	categoryEntity.Translations = make(map[string]string)

	for rows.Next() {
		var langCode, name string
		if err := rows.Scan(&langCode, &name); err != nil {
			return nil, fmt.Errorf("find one category by ID %v: %w", categoryID, err)
		}
		categoryEntity.Translations[langCode] = name
	}

	return &categoryEntity, nil
}

func (r *repo) FindByProjectID(c context.Context, projectID uuid.UUID) ([]*menu.Category, error) {
	var output []*menu.Category

	sql := "SELECT id, project_id, position FROM module_menu.categories WHERE project_id = $1"
	rows, err := r.db.Query(c, sql, projectID)
	if err != nil {
		return nil, fmt.Errorf("find one category by projectID %v: %w", projectID, err)
	}
	defer rows.Close()

	for rows.Next() {
		var cat menu.Category
		if err := rows.Scan(&cat.ID, &cat.ProjectID, &cat.Position); err != nil {
			return nil, fmt.Errorf("find one category by projectID %v: %w", projectID, err)
		}

		sql = "SELECT lang_code, name FROM module_menu.category_translations WHERE category_id = $1;"
		trRows, err := r.db.Query(c, sql, cat.ID)
		if err != nil {
			return nil, fmt.Errorf("find one category by projectID %v: %w", projectID, err)
		}

		cat.Translations = make(map[string]string)

		for trRows.Next() {
			var langCode, name string
			if err := trRows.Scan(&langCode, &name); err != nil {
				trRows.Close()
				return nil, fmt.Errorf("find one category by projectID %v: %w", projectID, err)
			}
			cat.Translations[langCode] = name
		}
		trRows.Close()

		output = append(output, &cat)
	}

	return output, nil
}

func (r *repo) Save(c context.Context, category *menu.Category) error {
	tx, err := r.db.Begin(c)
	if err != nil {
		return fmt.Errorf("save category: %w", err)
	}
	defer tx.Rollback(c)

	sql := `INSERT INTO module_menu.categories (id, project_id, position)
						VALUES ($1, $2, $3)
					ON CONFLICT(id)
					DO UPDATE SET position = EXCLUDED.position`
	_, err = tx.Exec(c, sql, category.ID, category.ProjectID, category.Position)
	if err != nil {
		return fmt.Errorf("save category: %w", err)
	}

	sql = "DELETE FROM module_menu.category_translations WHERE category_id = $1;"
	_, err = tx.Exec(c, sql, category.ID)
	if err != nil {
		return fmt.Errorf("save category: %w", err)
	}

	for key, value := range category.Translations {
		sql = "INSERT INTO module_menu.category_translations (category_id, lang_code, name) VALUES ($1, $2, $3);"
		_, err := tx.Exec(c, sql, category.ID, key, value)
		if err != nil {
			return fmt.Errorf("save category: %w", err)
		}
	}

	if err := tx.Commit(c); err != nil {
		return fmt.Errorf("save category: %w", err)
	}

	return nil
}

func (r *repo) Remove(c context.Context, category *menu.Category) error {
	tx, err := r.db.Begin(c)
	if err != nil {
		return fmt.Errorf("remove category: %w", err)
	}
	defer tx.Rollback(c)

	sql := "DELETE FROM module_menu.categories WHERE id = $1;"
	_, err = tx.Exec(c, sql, category.ID)
	if err != nil {
		return fmt.Errorf("remove category %v: %w", category.ID, err)
	}

	sql = "DELETE FROM module_menu.category_translations WHERE category_id = $1;"
	_, err = tx.Exec(c, sql, category.ID)
	if err != nil {
		return fmt.Errorf("remove category %v: %w", category.ID, err)
	}

	if err := tx.Commit(c); err != nil {
		return fmt.Errorf("remove category: %w", err)
	}

	return nil
}
