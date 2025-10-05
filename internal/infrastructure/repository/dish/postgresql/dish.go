package dish_postgresql_repository

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

func New(db *pgxpool.Pool) menu.DishRepository {
	return &repo{
		db: db,
	}
}

func (r *repo) FindOneByID(c context.Context, dishID uuid.UUID) (*menu.Dish, error) {
	var dishEntity menu.Dish

	sql := "SELECT id, category_id, position, photo_url, price, is_available FROM module_menu.dishes WHERE id = $1;"
	row := r.db.QueryRow(c, sql, dishID)
	err := row.Scan(
		&dishEntity.ID, &dishEntity.CategoryID, &dishEntity.Position,
		&dishEntity.PhotoURL, &dishEntity.Price, &dishEntity.IsAvailable,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("find one dish by ID %v: %w: %w", dishID, err, common.ErrNotFound)
		}
		return nil, fmt.Errorf("find one dish by ID %v: %w: %w", dishID, err, common.ErrInfrastructure)
	}

	sql = "SELECT lang, name, description FROM module_menu.dish_translations WHERE dish_id = $1"
	rows, err := r.db.Query(c, sql, dishID)
	if err != nil {
		return nil, fmt.Errorf("find one dish by ID %v: %w: %w", dishID, err, common.ErrInfrastructure)
	}
	defer rows.Close()

	dishEntity.Translations = make(map[string]menu.DishTranslation)

	for rows.Next() {
		var lang, name, description string
		if err := rows.Scan(&lang, &name, &description); err != nil {
			return nil, fmt.Errorf("find one dish by ID %v: %w: %w", dishID, err, common.ErrInfrastructure)
		}
		dishEntity.Translations[lang] = menu.DishTranslation{
			Name:        name,
			Description: description,
		}
	}

	options, err := r.findOptionGroups(c, dishID)
	if err != nil {
		return nil, fmt.Errorf("find one dish by ID %v: %w: %w", dishID, err, common.ErrInfrastructure)
	}
	dishEntity.Options = options

	return &dishEntity, nil
}

func (r *repo) FindByProjectID(c context.Context, projectID uuid.UUID) ([]*menu.Dish, error) {
	var output []*menu.Dish

	sql := "SELECT id, category_id, position, photo_url, price, is_available FROM module_menu.dishes WHERE project_id = $1;"
	rows, err := r.db.Query(c, sql, projectID)
	if err != nil {
		return nil, fmt.Errorf("find dishes by projectID %v: %w: %w", projectID, err, common.ErrInfrastructure)
	}

	for rows.Next() {
		var dish menu.Dish
		err := rows.Scan(
			&dish.ID, &dish.CategoryID, &dish.Position,
			&dish.PhotoURL, &dish.Price, &dish.IsAvailable,
		)
		if err != nil {
			return nil, fmt.Errorf("find dishes by projectID %v: %w: %w", projectID, err, common.ErrInfrastructure)
		}

		{
			sql = "SELECT lang, name, description FROM module_menu.dish_translations WHERE dish_id = $1"
			trRows, err := r.db.Query(c, sql, dish.ID)
			if err != nil {
				return nil, fmt.Errorf("find one dish by ID %v: %w: %w", dish.ID, err, common.ErrInfrastructure)
			}

			dish.Translations = make(map[string]menu.DishTranslation)

			for trRows.Next() {
				var lang, name, description string
				if err := trRows.Scan(&lang, &name, &description); err != nil {
					return nil, fmt.Errorf("find one dish by ID %v: %w: %w", dish.ID, err, common.ErrInfrastructure)
				}
				dish.Translations[lang] = menu.DishTranslation{
					Name:        name,
					Description: description,
				}
			}
		}

		options, err := r.findOptionGroups(c, dish.ID)
		if err != nil {
			return nil, fmt.Errorf("find dishes by projectID %v: %w: %w", projectID, err, common.ErrInfrastructure)
		}
		dish.Options = options

		output = append(output, &dish)
	}

	return output, nil
}

func (r *repo) Save(c context.Context, d *menu.Dish) error {
	tx, err := r.db.Begin(c)
	if err != nil {
		return fmt.Errorf("save dish %v: %w: %w", d.ID, err, common.ErrInfrastructure)
	}
	defer tx.Rollback(c)

	// dish
	sql := `INSERT INTO module_menu.dishes
						(id, category_id, position, photo_url, price, is_available)
						VALUES ($1, $2, $3, $4, $5, $6)
					ON CONFLICT (id)
					DO UPDATE SET
						category_id = EXCLUDED.category_id,
						position = EXCLUDED.position,
						photo_url = EXCLUDED.photo_url,
						price = EXCLUDED.price,
						is_available = EXCLUDED.is_available`
	_, err = tx.Exec(c, sql, d.ID, d.CategoryID, d.Position, d.PhotoURL, d.Price, d.IsAvailable)
	if err != nil {
		return fmt.Errorf("save dish %v: %w: %w", d.ID, err, common.ErrInfrastructure)
	}

	// translation
	_, err = tx.Exec(c, "DELETE FROM module_menu.dish_translations WHERE dish_id = $1", d.ID)
	if err != nil {
		return fmt.Errorf("save dish %v: %w: %w", d.ID, err, common.ErrInfrastructure)
	}

	for lang, tr := range d.Translations {
		sql := `INSERT INTO module_menu.dish_translation 
					(dish_id, lang_code, name, description)
					VALUES ($1, $2, $3, $4);`
		_, err := tx.Exec(c, sql, d.ID, lang, tr.Name, tr.Description)
		if err != nil {
			return fmt.Errorf("save dish %v: %w: %w", d.ID, err, common.ErrInfrastructure)
		}
	}

	if err := r.addOptions(c, tx, d.Options); err != nil {
		return fmt.Errorf("save dish %v: %w: %w", d.ID, err, common.ErrInfrastructure)
	}

	return nil
}

func (r *repo) Remove(c context.Context, dish *menu.Dish) error {
	sql := "DELETE FROM module_menu.dishes WHERE id = $1;"

	_, err := r.db.Exec(c, sql, dish.ID)
	if err != nil {
		return fmt.Errorf("remove dish %v: %w: %w", dish.ID, err, common.ErrInfrastructure)
	}

	return nil
}

func (r *repo) findOptionGroups(c context.Context, dishID uuid.UUID) ([]menu.OptionGroup, error) {
	var output []menu.OptionGroup

	sql := "SELECT id, dish_id, position, min_quantity, max_quantity FROM module_menu.option_groups WHERE dish_id = $1;"
	rows, err := r.db.Query(c, sql, dishID)
	if err != nil {
		return nil, fmt.Errorf("find one dish by ID %v: %w: %w", dishID, err, common.ErrInfrastructure)
	}

	for rows.Next() {
		var opt menu.OptionGroup
		if err := rows.Scan(&opt.ID, &opt.DishID, &opt.Position, &opt.MinQuantity, &opt.MaxQuantity); err != nil {
			return nil, fmt.Errorf("find one dish by ID %v: %w: %w", dishID, err, common.ErrInfrastructure)
		}

		{
			sql = "SELECT lang, name FROM module_menu.option_group_translations WHERE option_group_id = $1;"
			trRows, err := r.db.Query(c, sql, opt.ID)
			if err != nil {
				return nil, fmt.Errorf("find one dish by ID %v: %w: %w", dishID, err, common.ErrInfrastructure)
			}

			for trRows.Next() {
				var lang, name string
				if err := trRows.Scan(&lang, &name); err != nil {
					trRows.Close()
					return nil, fmt.Errorf("find one dish by ID %v: %w: %w", dishID, err, common.ErrInfrastructure)
				}
				opt.Translations[lang] = name
			}
		}

		items, err := r.findOptionItems(c, opt.ID)
		if err != nil {
			return nil, fmt.Errorf("find one dish by ID %v: %w: %w", dishID, err, common.ErrInfrastructure)
		}
		opt.Items = items

		output = append(output, opt)
	}

	return output, nil
}

func (r *repo) findOptionItems(c context.Context, optionID uuid.UUID) ([]menu.OptionItem, error) {
	var output []menu.OptionItem

	sql := `SELECT id, option_group_id, position, min_quantity, max_quantity, price_modifier
					FROM module_menu.option_items
					WHERE option_group_id = $1;`
	rows, err := r.db.Query(c, sql, optionID)
	if err != nil {
		return nil, fmt.Errorf("find items by option ID %v: %w: %w", optionID, err, common.ErrInfrastructure)
	}

	for rows.Next() {
		var item menu.OptionItem
		err := rows.Scan(&item.ID, &item.OptionGroupID, &item.Position, &item.MinQuantity, &item.MaxQuantity, &item.PriceModifier)
		if err != nil {
			return nil, fmt.Errorf("find items by option ID %v: %w: %w", optionID, err, common.ErrInfrastructure)
		}

		sql := "SELECT lang, name FROM module_menu.option_item_translations WHERE item_id = $1;"
		trRows, err := r.db.Query(c, sql, item.ID)
		if err != nil {
			return nil, fmt.Errorf("find items by option ID %v: %w: %w", optionID, err, common.ErrInfrastructure)
		}

		for trRows.Next() {
			var lang, name string
			if err := trRows.Scan(&lang, &name); err != nil {
				return nil, fmt.Errorf("find items by option ID %v: %w: %w", optionID, err, common.ErrInfrastructure)
			}
			item.Translations[lang] = name
		}

		output = append(output, item)
	}

	return output, nil
}

func (r *repo) addOptions(c context.Context, tx pgx.Tx, options []menu.OptionGroup) error {
	sql := `INSERT INTO module_menu.option_groups
						(id, dish_id, position, min_quantity, max_quantity)
						VALUES ($1, $2, $3, $4, $5)
					ON CONFLICT (id)
					DO UPDATE SET
						position = EXCLUDED.position,
						min_quantity = EXCLUDED.min_quantity,
						max_quantity = EXCLUDED.max_quantity;`

	for _, opt := range options {
		_, err := tx.Exec(c, sql, opt.ID, opt.DishID, opt.Position, opt.MinQuantity, opt.MaxQuantity)
		if err != nil {
			return fmt.Errorf("add options for dish %v: %w", opt.DishID, err)
		}

		for lang, name := range opt.Translations {
			sql := `INSERT INTO module_menu.option_group_translations
								(option_group_id, lang_code, name)
								VALUES ($1, $2, $3)
							ON CONFLICT (option_group_id, lang_code)
							DO UPDATE SET
								name = EXCLUDED.name;`
			_, err := tx.Exec(c, sql, opt.ID, lang, name)
			if err != nil {
				return fmt.Errorf("add options for dish %v: %w", opt.DishID, err)
			}
		}

		if err := r.addItems(c, tx, opt.Items); err != nil {
			return fmt.Errorf("add options for dish %v: %w", opt.DishID, err)
		}
	}

	return nil
}

func (r *repo) addItems(c context.Context, tx pgx.Tx, items []menu.OptionItem) error {
	sql := `INSERT INTO module_menu.option_items
						(id, option_group_id, position, min_quantity, max_quantity, price_modifier)
						VALUES ($1, $2, $3, $4, $5, $6)
					ON CONFLICT (id)
					DO UPDATE SET
						position = EXCLUDED.position,
						min_quantity = EXCLUDED.min_quantity,
						max_quantity = EXCLUDED.max_quantity,
						price_modifier = EXCLUDED.price_modifier;`

	sql2 := `INSERT INTO module_menu.option_item_translations
							(item_id, lang_code, name)
							VALUES ($1, $2, $3)
						ON CONFLICT (item_id, lang_code),
						DO UPDATE SET`

	for _, item := range items {
		_, err := tx.Exec(c, sql, item.ID, item.OptionGroupID, item.Position, item.MinQuantity, item.MaxQuantity, item.PriceModifier)
		if err != nil {
			return fmt.Errorf("add items for option %v: %w", item.OptionGroupID, err)
		}

		for lang, name := range item.Translations {
			_, err := tx.Exec(c, sql2, item.ID, lang, name)
			if err != nil {
				return fmt.Errorf("add items for option %v: %w", item.OptionGroupID, err)
			}
		}
	}
	return nil
}
