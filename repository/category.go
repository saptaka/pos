package repository

import (
	"context"
	"database/sql"

	"github.com/saptaka/pos/model"
)

type CategoryRepo interface {
	GetCategoryByID(ctx context.Context, id int64) (model.Category, error)
	GetCategories(ctx context.Context, limit, skip int) ([]model.Category, error)
	UpdateCategory(ctx context.Context, category model.Category) error
	CreateCategory(ctx context.Context, name string) (model.Category, error)
	DeleteCategory(ctx context.Context, id int64) error
}

func (r repo) GetCategoryByID(ctx context.Context, id int64) (model.Category, error) {
	var category model.Category
	query := "SELECT id, name FROM categories WHERE id=?"
	rows := r.db.QueryRowContext(ctx, query, id)
	err := rows.Scan(&category.CategoryId, &category.Name)
	if err != nil {
		return category, err
	}

	return category, nil
}

func (r repo) GetCategories(ctx context.Context,
	limit, skip int) ([]model.Category, error) {
	query := "SELECT id, name FROM categories "
	var rows *sql.Rows
	var err error
	if limit > 0 {
		query += " limit ? offset ?;"
		rows, err = r.db.QueryContext(ctx, query, limit, skip)
		if err != nil {
			return nil, err
		}
	} else {
		rows, err = r.db.QueryContext(ctx, query)
		if err != nil {
			return nil, err
		}
	}
	var categories []model.Category
	for rows.Next() {
		var category model.Category
		err := rows.Scan(&category.CategoryId, &category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (r repo) UpdateCategory(ctx context.Context,
	category model.Category) error {
	_, err := r.GetCategoryByID(ctx, category.CategoryId)
	if err != nil {
		return err
	}
	query := "UPDATE categories SET name=? WHERE id=?"
	_, err = r.db.ExecContext(ctx, query, category.Name,
		category.CategoryId)
	if err != nil {
		return err
	}

	return err
}

func (r repo) CreateCategory(ctx context.Context, name string) (model.Category, error) {
	var categoryDetail model.Category
	insertQuery := `INSERT INTO 
		categories (name) 
	VALUES (?);`
	stmt, err := r.db.PrepareContext(ctx, insertQuery)
	if err != nil {
		return categoryDetail, err
	}
	res, err := stmt.Exec(name)
	if err != nil {
		return categoryDetail, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return categoryDetail, err
	}
	selectQuery := `SELECT id, 
					name,
					updated_at, 
					created_at
					FROM categories 
					WHERE id=?`
	rows := r.db.QueryRowContext(ctx, selectQuery, id)

	err = rows.Scan(
		&categoryDetail.CategoryId,
		&categoryDetail.Name,
		&categoryDetail.UpdatedAt,
		&categoryDetail.CreatedAt)
	return categoryDetail, err

}

func (r repo) DeleteCategory(ctx context.Context, id int64) error {
	_, err := r.GetCategoryByID(ctx, id)
	if err != nil {
		return err
	}
	query := "DELETE FROM categories WHERE id=?"
	_, err = r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return err
}
