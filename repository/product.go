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
	GetProductByID(ctx context.Context, id int64) (model.Product, error)
	GetProducts(ctx context.Context, limit, skip int, product model.Product) ([]model.Product, error)
	UpdateProduct(ctx context.Context, product model.Product) error
	CreateProduct(ctx context.Context, product model.ProductCreateRequest) (model.Product, error)
	DeleteProduct(ctx context.Context, id int64) error
	GetProductsByIds(ctx context.Context, ids []int64) ([]model.Product, error)
}

func (r repo) GetProductByID(ctx context.Context, id int64) (model.Product, error) {
	var product model.Product
	var category model.Category
	query := `SELECT 
				id,
				name,
				stock,
				price,
				image,
				category_id,
				sku,
				discount_id
			FROM products 
			WHERE id=?`
	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&product.ProductId,
		&product.Name,
		&product.Stock,
		&product.Price,
		&product.Image,
		&product.CategoryId,
		&product.SKU,
		&product.DiscountId,
	)
	if err != nil {
		return product, err
	}

	var discountById *model.Discount
	if product.DiscountId != nil {
		discount, err := r.GetDiscountByID(ctx, *product.DiscountId)
		if err != sql.ErrNoRows && err != nil {
			log.Println("error get product ", err)
			return product, err
		}

		if discount.Type == model.BuyN {
			discount.StringFormat = fmt.Sprintf("Buy %d only Rp. %s",
				discount.Qty, fmt.Sprintf("%d", discount.Result))
		} else {
			discountResult := fmt.Sprint(discount.Result, "%")
			discountPrice := product.Price - (product.Price * discount.Result / 100)
			discount.StringFormat = fmt.Sprintf("Discount %s Rp. %s",
				discountResult, fmt.Sprintf("%d", discountPrice))
		}

		discountById = &discount
	}

	product.Discount = discountById

	if product.CategoryId != nil {
		category, err = r.GetCategoryByID(ctx, *product.CategoryId)
		if err != nil {
			return product, err
		}
		if product.DiscountId != nil {
			discount, err := r.GetDiscountByID(ctx, *product.DiscountId)
			if err != nil {
				return product, err
			}
			product.Discount = &discount
		}

		product.Category = &category
	}

	return product, nil
}

func (r repo) GetProducts(ctx context.Context,
	limit, skip int, product model.Product) ([]model.Product, error) {

	productChan := make(chan []model.Product)
	go func(productChanData chan []model.Product) {
		querySelect := `SELECT 
				id,
				name,
				stock,
				price,
				image,
				category_id ,
				sku,
				discount_id
			FROM products 
			%s 
			`
		var withQuery string
		values := make([]interface{}, 0)
		if product.Name != "" {
			withQuery = " WHERE name LIKE CONCAT('%',?,'%')"
			values = append(values, product.Name)
		} else if product.CategoryId != nil {
			withQuery = " WHERE category_id=?"
			values = append(values, *product.CategoryId)
		}
		querySelect = fmt.Sprintf(querySelect, withQuery)

		var rows *sql.Rows
		var err error
		if limit > 0 {
			values = append(values, limit, skip)
			querySelect += " limit ? offset ?;"
			rows, err = r.db.QueryContext(ctx, querySelect, values...)
			if err != nil {
				log.Println("error get product ", err)
				productChanData <- make([]model.Product, 0)
				return
			}
		} else {
			if len(values) > 0 {
				rows, err = r.db.QueryContext(ctx, querySelect, values...)
			} else {
				rows, err = r.db.QueryContext(ctx, querySelect)
			}
			if err != nil {
				log.Println("error get product ", err)
				productChanData <- make([]model.Product, 0)
				return
			}
		}
		var products []model.Product
		for rows.Next() {
			var product model.Product
			err := rows.Scan(
				&product.ProductId,
				&product.Name,
				&product.Stock,
				&product.Price,
				&product.Image,
				&product.CategoryId,
				&product.SKU,
				&product.DiscountId,
			)
			if err != nil {
				log.Println("error get product ", err)
				productChanData <- make([]model.Product, 0)
				return
			}
			if product.DiscountId != nil {
				discount, err := r.GetDiscountByID(ctx, *product.DiscountId)
				if err != sql.ErrNoRows && err != nil {
					log.Println("error get product ", err)
					return
				}

				if discount.Type == model.BuyN {
					discount.StringFormat = fmt.Sprintf("Buy %d only Rp. %s",
						discount.Qty, fmt.Sprintf("%d", discount.Result))
				} else {
					discountResult := fmt.Sprint(discount.Result, "%")
					discountPrice := product.Price - (product.Price * discount.Result / 100)
					discount.StringFormat = fmt.Sprintf("Discount %s Rp. %s",
						discountResult, fmt.Sprintf("%d", discountPrice))
				}

				product.Discount = &discount
			}

			products = append(products, product)
		}
		productChanData <- products

	}(productChan)

	categoryChan := make(chan []model.Category)
	go func(categoryChanData chan []model.Category) {
		categories, err := r.GetCategories(ctx, 0, 0)
		if err != nil {
			log.Println("error get categories ", err)
			categoryChanData <- make([]model.Category, 0)
			return
		}
		categoryChanData <- categories
	}(categoryChan)

	products := <-productChan
	categories := <-categoryChan
	categoriesMap := make(map[int64]*model.Category)
	for _, category := range categories {
		categoriesMap[category.CategoryId] = &category
		categories = categories[1:]
	}
	for index, product := range products {
		if product.CategoryId != nil {
			products[index].Category = categoriesMap[*product.CategoryId]
		}
	}

	if products == nil {
		products = make([]model.Product, 0)
	}

	return products, nil
}

func (r repo) UpdateProduct(ctx context.Context,
	Product model.Product) error {

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
	if Product.CategoryId != nil && *Product.CategoryId != 0 {
		query += " category_id=?,"
		values = append(values, Product.CategoryId)
	}

	if countUpdate > 0 {
		query += " updated_at=CURRENT_TIMESTAMP()  WHERE id=? "
		values = append(values, Product.ProductId)
	} else {
		return fmt.Errorf("nothing updated")
	}

	_, err := r.db.ExecContext(ctx, query, values...)
	if err != nil {
		return err
	}

	return nil
}

func (r repo) CreateProduct(ctx context.Context, product model.ProductCreateRequest) (model.Product, error) {

	var productDetail model.Product

	insertQuery := `INSERT INTO 
		products (name,image, price, stock, category_id,
			 updated_at, created_at) 
	VALUES (?,?,?,?,?,?,?);`

	stmt, err := r.db.PrepareContext(ctx, insertQuery)
	if err != nil {
		return productDetail, err
	}

	now, _ := time.Parse(model.RFC3339MilliZ, time.Now().UTC().Format(model.RFC3339MilliZ))

	result, err := stmt.Exec(
		product.Name,
		product.Image,
		product.Price,
		product.Stock,
		product.CategoryId,
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

	go func(repoInside repo, id int64, discount *model.Discount) {
		var discountId *int64
		if discount != nil {
			discountIDResult, err := repoInside.CreateDiscount(ctx, *discount)
			if err != nil {
				log.Println(err)
				return
			}
			discountId = &discountIDResult
		}

		updateQuery := `UPDATE products 
				SET discount_id=?, sku=?
				WHERE id=?`

		_, err = repoInside.db.ExecContext(ctx, updateQuery, discountId, fmt.Sprintf("ID%03d", id), id)
		if err != nil {
			log.Println(err)
			return
		}
	}(r, id, product.Discount)

	productDetail = model.Product{

		ProductId: id,
		Name:      product.Name,
		Stock:     product.Stock,
		SKU:       fmt.Sprintf("ID%03d", id),
		Price:     product.Price,
		Image:     product.Image,
		UpdatedAt: &now,
		CreatedAt: &now,

		CategoryId: product.CategoryId,
	}

	return productDetail, err

}

func (r repo) DeleteProduct(ctx context.Context, id int64) error {
	query := "DELETE FROM products WHERE id=?"
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return err
}

func (r repo) GetProductsByIds(ctx context.Context, ids []int64) ([]model.Product, error) {
	if len(ids) == 0 {
		return nil, sql.ErrNoRows
	}
	querySelect := `SELECT 
				id,
				name,
				price,
				discount_id  
			FROM products 
			WHERE id IN (%s) ORDER BY id ASC`
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

	if err != nil {
		return nil, err
	}
	var products []model.Product
	for rows.Next() {
		var product model.Product
		err := rows.Scan(
			&product.ProductId,
			&product.Name,
			&product.Price,
			&product.DiscountId,
		)
		if err != nil {
			return products, err
		}
		if product.DiscountId != nil {
			discount, err := r.GetDiscountByID(ctx, *product.DiscountId)
			if err != nil {
				return products, err
			}
			product.Discount = &discount
		}

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

func (r repo) GetDiscountByID(ctx context.Context, id int64) (model.Discount, error) {
	var discount model.Discount
	query := `SELECT 
			id,
			qty,
			types,
			result,
			expired_at,
			expired_at_format,
			string_format
			FROM discounts  
			WHERE id=?
			`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&discount.DiscountID,
		&discount.Qty,
		&discount.Type,
		&discount.Result,
		&discount.ExpiratedAt,
		&discount.ExpiredAtFormat,
		&discount.StringFormat,
	)
	return discount, err
}
