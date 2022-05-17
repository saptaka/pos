package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

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
			`

	values := make([]interface{}, 0)
	if query != "" {
		withQuery = " WHERE p.name=?"
		values = append(values, query)
	}
	querySelect = fmt.Sprintf(querySelect, withQuery)

	var rows *sql.Rows
	var err error
	if limit > 0 {
		values = append(values, limit, skip)
		querySelect += " limit ? offset ?;"
		rows, err = r.db.QueryContext(ctx, querySelect, values...)
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		if err != nil {
			return nil, err
		}
	} else {
		if len(values) > 0 {
			rows, err = r.db.QueryContext(ctx, querySelect, values...)
		} else {
			rows, err = r.db.QueryContext(ctx, querySelect)
		}
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		if err != nil {
			return nil, err
		}
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
	_, err := r.GetProductByID(ctx, int(Product.ProductId))
	if err == sql.ErrNoRows {
		return sql.ErrNoRows
	}

	if err != nil {
		return err
	}

	query := `UPDATE products SET `
	var countUpdate int
	var values []interface{}
	if Product.Name != "" {
		countUpdate++
		query += " name=?,"
		values = append(values, Product.Name)
	}
	if Product.Image != "" {
		query += " 	image=?,"
		values = append(values, Product.Image)
	}
	if Product.Stock != 0 {
		query += " 	stock=?,"
		values = append(values, Product.Stock)
	}
	if Product.Price != 0 {
		query += " price=?,"
		values = append(values, Product.Price)
	}
	if Product.CategoryID != 0 {
		query += " category_id=?,"
		values = append(values, Product.CategoryID)
	}

	if countUpdate > 0 {
		query += " updated_at=CURRENT_TIMESTAMP()  WHERE id=? "
		values = append(values, Product.ProductId)
	} else {
		return fmt.Errorf("nothing updated")
	}

	result, err := r.db.ExecContext(ctx, query, values...)
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

	insertQuery := `INSERT INTO 
		products (name,image, price, stock, 
			 updated_at, created_at) 
	VALUES (?,?,?,?,?,?);`

	stmt, err := r.db.PrepareContext(ctx, insertQuery)
	if err != nil {
		return productDetail, err
	}

	now := time.Now()

	result, err := stmt.Exec(
		product.Name,
		product.Image,
		product.Price,
		product.Stock,
		now,
		now,
	)
	if err != nil {
		return productDetail, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return productDetail, err
	}

	go func(repoInside repo, id, categoryId int64, discount *model.Discount) {
		var discountId int64
		if discount != nil {
			discountIDResult, err := repoInside.CreateDiscount(ctx, *discount)
			if err != nil {
				log.Println(err)
				return
			}
			discountId = discountIDResult
		}
		updateQuery := `UPDATE products 
				SET sku=CONCAT('ID',LPAD(?,3,0)), discount_id=?, category_id=?
				WHERE id=?`
		_, err = repoInside.db.ExecContext(ctx, updateQuery, id, discountId, categoryId, id)
		if err != nil {
			log.Println(err)
			return
		}
	}(r, id, product.CategoryID, product.Discount)

	productDetail = model.Product{
		ProductId:  id,
		Name:       product.Name,
		SKU:        "ID" + fmt.Sprintf("|%03d|", id),
		Stock:      product.Stock,
		Price:      product.Price,
		Image:      product.Image,
		CategoryID: product.CategoryID,
		UpdatedAt:  &now,
		CreatedAt:  &now,
	}

	return productDetail, err

}

func (r repo) DeleteProduct(ctx context.Context, id int) error {
	_, err := r.GetProductByID(ctx, id)
	if err == sql.ErrNoRows {
		return sql.ErrNoRows
	}

	if err != nil {
		return err
	}
	query := "DELETE FROM products WHERE id=?"
	result, err := r.db.ExecContext(ctx, query, id)
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
		return nil, sql.ErrNoRows
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
