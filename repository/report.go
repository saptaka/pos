package repository

import (
	"context"

	"github.com/saptaka/pos/model"
)

type ReportRepo interface {
	GetRevenues(ctx context.Context) (model.Revenue, error)
	GetSolds(ctx context.Context) (model.Solds, error)
}

func (r repo) GetRevenues(ctx context.Context) (model.Revenue, error) {
	query := `
		SELECT
		payments.id,
		payments.logo,
		payments.name,
		payments.types,
		orders.total_paid
	FROM
		payments
		JOIN orders
	WHERE
		payments.id = orders.payment_type_id
	`
	var revenue model.Revenue
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return revenue, err
	}

	var totalRevenue int
	for rows.Next() {
		var payment model.PaymentTypeItem
		err := rows.Scan(
			&payment.PaymentId,
			&payment.Logo,
			&payment.Name,
			&payment.Type,
			&payment.TotalAmount,
		)
		if err != nil {
			return revenue, nil
		}
		totalRevenue += payment.TotalAmount
		revenue.PaymentType = append(revenue.PaymentType, payment)
	}
	revenue.TotalRevenue = totalRevenue
	return revenue, nil
}

func (r repo) GetSolds(ctx context.Context) (model.Solds, error) {
	query := `
		SELECT
		ordered_products.product_id,
		products.name,
		SUM(ordered_products.qty) as totalAQty,
		SUM(total_normal_price) as totalAmount
	FROM
		ordered_products
		JOIN products ON ordered_products.product_id = products.id
		GROUP BY ordered_products.product_id
	`
	var sold model.Solds
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return sold, err
	}

	for rows.Next() {
		var soldProduct model.SoldProduct
		err := rows.Scan(
			&soldProduct.ProductId,
			&soldProduct.Name,
			&soldProduct.TotalQty,
			&soldProduct.TotalAmount,
		)
		if err != nil {
			return sold, nil
		}
		sold.OrderProduct = append(sold.OrderProduct, soldProduct)
	}
	return sold, nil
}
