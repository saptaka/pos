package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/saptaka/pos/model"
)

type ProductRepo interface {
	GetProductByID(ctx context.Context, id int) (model.Product, error)
	GetProducts(ctx context.Context, limit, skip int, query string) ([]model.Product, error)
	UpdateProduct(ctx context.Context, product model.Product) error
	CreateProduct(ctx context.Context, product model.Product) (model.Product, error)
	DeleteProduct(ctx context.Context, id int) error
	GetProductsByIds(ctx context.Context, ids []int64) ([]model.Product, error)
}

func (r repo) GetProductByID(ctx context.Context, id int) (model.Product, error) {
	var product model.Product
	var category model.Category
	query := `SELECT 
				p.id,
				p.name,
				p.sku,
				p.stock,
				p.price,
				p.image,
				c.id,
				c.name 
			FROM products p JOIN categories c
			ON p.category_id=c.id
			WHERE p.id=?`
	rows := r.db.QueryRowContext(ctx, query, id)
	err := rows.Scan(
		&product.ProductId,
		&product.Name,
		&product.SKU,
		&product.Stock,
		&product.Price,
		&product.Image,
		&category.CategoryId,
		&category.Name,
	)
	product.Category = &category
	if err != nil {
		return product, err
	}

	return product, nil
}

func (r repo) GetProducts(ctx context.Context,
	limit, skip int, query string) ([]model.Product, error) {
	var withQuery string

	querySelect := `SELECT 
				p.id,
				p.name,
				p.sku,
				p.stock,
				p.price,
				p.image,
				c.id,
				c.name 
			FROM products p JOIN categories c 
			ON p.category_id=c.id 
			%s 
			limit ? offset ?`
	values := make([]interface{}, 0)
	if query != "" {
		withQuery = " WHERE p.name=?"
		values = append(values, query)
	}
	querySelect = fmt.Sprintf(querySelect, withQuery)
	values = append(values, limit, skip)
	rows, err := r.db.QueryContext(ctx, querySelect, values...)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var products []model.Product
	for rows.Next() {
		var product model.Product
		var category model.Category
		err := rows.Scan(
			&product.ProductId,
			&product.Name,
			&product.SKU,
			&product.Stock,
			&product.Price,
			&product.Image,
			&category.CategoryId,
			&category.Name,
		)
		product.Category = &category
		if err != nil {
			return products, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (r repo) UpdateProduct(ctx context.Context,
	Product model.Product) error {
	query := `UPDATE products SET 
			name=?, 
			image=?,
			stock=?,
			price=?,
			category_id=?,
			updated_at=CURRENT_TIMESTAMP() WHERE id=?`
	result, err := r.db.ExecContext(ctx, query,
		Product.Name,
		Product.Image,
		Product.Stock,
		Product.Price,
		Product.CategoryID,
		Product.ProductId)
	if err != nil {
		return err
	}
	rowAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r repo) CreateProduct(ctx context.Context, product model.Product) (model.Product, error) {

	var productDetail model.Product
	var discountId *int64
	if product.Discount != nil {
		discountIDResult, err := r.CreateDiscount(ctx, *product.Discount)
		if err != nil {
			return productDetail, err
		}
		discountId = &discountIDResult
	}
	insertQuery := `INSERT INTO 
		products (name,image, price, stock, category_id, discount_id) 
	VALUES (?,?,?,?,?,?);`

	stmt, err := r.db.PrepareContext(ctx, insertQuery)
	if err != nil {
		return productDetail, err
	}

	result, err := stmt.Exec(
		product.Name,
		product.Image,
		product.Price,
		product.Stock,
		product.CategoryID,
		discountId,
	)
	if err != nil {
		return productDetail, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return productDetail, err
	}

	updateQuery := `UPDATE products 
				SET sku=CONCAT('ID',LPAD(?,3,0))
				WHERE id=?`
	_, err = r.db.ExecContext(ctx, updateQuery, id, id)
	if err != nil {
		return productDetail, err
	}

	selectQuery := `SELECT id, 
					name,
					category_id,
					sku,
					stock,
					price,
					image,
					updated_at, 
					created_at
					FROM products 
					WHERE id=?`
	rows := r.db.QueryRowContext(ctx, selectQuery, id)
	err = rows.Scan(
		&productDetail.ProductId,
		&productDetail.Name,
		&productDetail.CategoryID,
		&productDetail.SKU,
		&productDetail.Stock,
		&productDetail.Price,
		&productDetail.Image,
		&productDetail.UpdatedAt,
		&productDetail.CreatedAt)
	return productDetail, err

}

func (r repo) DeleteProduct(ctx context.Context, id int) error {
	query := "DELETE FROM products WHERE id=?"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r repo) GetProductsByIds(ctx context.Context, ids []int64) ([]model.Product, error) {
	if len(ids) == 0 {
		return nil, sql.ErrNoRows
	}
	querySelect := `SELECT 
				p.id,
				p.name,
				p.price,
				d.id,
				d.qty,
				d.types,
				d.result,
				d.expired_at,
				d.expired_at_format,
				d.string_format
			FROM products p JOIN discounts d 
			ON p.discount_id=d.id 
			WHERE p.id IN (%s) ORDER BY p.id ASC`
	values := make([]interface{}, 0)
	for _, id := range ids {
		values = append(values, id)
	}
	template := "?"
	if len(ids) > 1 {
		template += strings.Repeat(",?", len(ids)-1)
	}

	querySelect = fmt.Sprintf(querySelect, template)
	rows, err := r.db.QueryContext(ctx, querySelect, values...)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var products []model.Product
	for rows.Next() {
		var product model.Product
		var discount model.Discount
		err := rows.Scan(
			&product.ProductId,
			&product.Name,
			&product.Price,
			&discount.DiscountID,
			&discount.Qty,
			&discount.Type,
			&discount.Result,
			&discount.ExpiratedAt,
			&discount.ExpiredAtFormat,
			&discount.StringFormat,
		)
		if err != nil {
			return products, err
		}
		product.Discount = &discount
		products = append(products, product)
	}
	return products, nil
}

func (r repo) CreateDiscount(ctx context.Context, discount model.Discount) (int64, error) {

	query := `INSERT INTO 
	discounts (
		qty,
		types,
		result, 
		expired_at,
		expired_at_format, 
		string_format) 
	VALUES (?,?,?,FROM_UNIXTIME(?),FROM_UNIXTIME(?, '%d %M %Y'),?);`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(
		discount.Qty,
		discount.Type,
		discount.Result,
		discount.ExpiratedAt,
		discount.ExpiratedAt,
		discount.StringFormat,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return id, err
	}

	return id, err
}
