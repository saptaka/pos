package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/saptaka/pos/model"
)

type OrderRepo interface {
	GetOrder(ctx context.Context, limit, skip int) ([]model.Order, error)
	GetOrderByID(ctx context.Context, id int64) (model.Order, error)
	UpdateOrder() error
	CreateOrder(ctx context.Context,
		orderRequest model.Order) (model.Order, error)
	DownloadReceipt(ctx context.Context, id int64) (string, error)
	GetDownloadStatus(ctx context.Context, id int64) (bool, error)
}

func (r repo) GetOrder(ctx context.Context, limit, skip int) ([]model.Order, error) {

	querySelect := `
			SELECT
			o.id,
			o.payment_type_id,
			o.cashier_id,
			o.total_price,
			o.total_paid,
			o.total_return,
			o.receipt_id,
			o.created_at,
			COALESCE(c.name, '') as name,
			p.logo,
			p.name,
			p.types
		FROM
			orders o
			LEFT JOIN cashiers c ON o.cashier_id = c.id
			LEFT JOIN payments p ON o.payment_type_id = p.id
			LIMIT ? OFFSET ?;`

	rows, err := r.db.QueryContext(ctx, querySelect, limit, skip)
	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		return nil, err
	}
	var orders []model.Order
	for rows.Next() {
		var cashier model.Cashier
		var payment model.Payment
		var order model.Order
		err := rows.Scan(
			&order.OrderId,
			&order.PaymentID,
			&order.CashierID,
			&order.TotalPrice,
			&order.TotalPaid,
			&order.TotalReturn,
			&order.ReceiptID,
			&order.CreatedAt,
			&cashier.Name,
			&payment.Logo,
			&payment.Name,
			&payment.Type,
		)
		if err != nil {
			return orders, err
		}
		cashier.ChashierId = *order.CashierID
		payment.PaymentId = *order.PaymentID
		order.Cashier = cashier
		order.PaymentType = payment
		orders = append(orders, order)
	}
	return orders, nil
}

func (r repo) GetOrderByID(ctx context.Context, id int64) (model.Order, error) {

	querySelect := `
		SELECT 
		id,
		payment_type_id,
		cashier_id,
		total_price,
		total_paid,
		total_return,
		receipt_id,
		created_at
		FROM 
		orders 
		WHERE id=?;`

	rows := r.db.QueryRowContext(ctx, querySelect, id)
	var order model.Order
	err := rows.Scan(
		&order.OrderId,
		&order.PaymentID,
		&order.CashierID,
		&order.TotalPrice,
		&order.TotalPaid,
		&order.TotalReturn,
		&order.ReceiptID,
		&order.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return order, sql.ErrNoRows
	}
	if err != nil {
		return order, err
	}

	cashierChan := make(chan model.Cashier)
	if order.CashierID != nil {
		go func(id int64, cashierChanData chan model.Cashier) {
			querySelect := `
		SELECT
		id,name 
		FROM 
		cashier  
		WHERE id=?;`

			rows := r.db.QueryRowContext(ctx, querySelect, id)

			var cashier model.Cashier
			err := rows.Scan(
				&cashier.ChashierId,
				&cashier.Name,
			)
			if err == sql.ErrNoRows {
				cashierChanData <- model.Cashier{}
				log.Println("no cashier found")
				return
			}
			if err != nil {
				cashierChanData <- model.Cashier{}
				log.Println("error selecting cashier found")
				return
			}
			cashierChanData <- cashier

		}(*order.CashierID, cashierChan)
	} else {
		close(cashierChan)
	}

	paymentChan := make(chan model.Payment)

	if order.PaymentID != nil {

		go func(id int64, paymentChanData chan model.Payment) {
			querySelect := `
				SELECT
				id,name, types, logo
				FROM 
				payments  
				WHERE id=?;`

			rows := r.db.QueryRowContext(ctx, querySelect, id)

			var payment model.Payment
			err := rows.Scan(
				&payment.PaymentId,
				&payment.Name,
				&payment.Type,
				&payment.Logo,
			)
			if err == sql.ErrNoRows {
				paymentChanData <- model.Payment{}
				log.Println("no payment found")
				return
			}
			if err != nil {
				paymentChanData <- model.Payment{}
				log.Println("error selecting payment found : ")
				return
			}
			paymentChanData <- payment
		}(*order.PaymentID, paymentChan)
	} else {
		close(cashierChan)
	}
	order.Cashier = <-cashierChan
	order.PaymentType = <-paymentChan

	return order, nil
}

func (r repo) UpdateOrder() error {
	return nil
}

func (r repo) CreateOrder(ctx context.Context, orderRequest model.Order) (model.Order, error) {

	query := `INSERT INTO orders(payment_type_id, total_price, total_paid, total_return, created_at)
			VALUES (?,?,?,?,?);`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return orderRequest, err
	}
	res, err := stmt.Exec(
		orderRequest.PaymentID,
		orderRequest.TotalPrice,
		orderRequest.TotalPaid,
		orderRequest.TotalReturn,
		orderRequest.CreatedAt,
	)
	if err != nil {
		return orderRequest, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return orderRequest, err
	}
	orderRequest.OrderId = id
	err = r.CreateOrderedProduct(ctx, id, orderRequest.OrderedProduct)
	if err != nil {
		return orderRequest, err
	}
	return orderRequest, nil
}

func (r repo) DownloadReceipt(ctx context.Context, id int64) (string, error) {
	order, err := r.GetOrderByID(ctx, id)
	if err != nil {
		return "", err
	}
	err = r.updateIsDownloadReceipt(ctx, id)
	if err != nil {
		return "", err
	}
	return order.ReceiptIDFilePath, nil
}

func (r repo) updateIsDownloadReceipt(ctx context.Context, id int64) error {
	query := `
		UPDATE orders SET is_downloaded = 1 WHERE id=?
	`
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
	return nil
}

func (r repo) GetDownloadStatus(ctx context.Context, id int64) (bool, error) {
	query := `
		SELECT is_downloaded FROM orders WHERE id=?
	`
	var isDownload int
	rows := r.db.QueryRowContext(ctx, query, id)
	err := rows.Scan(&isDownload)
	if err != nil {
		return false, err
	}
	if isDownload == 0 {
		return false, nil
	}

	return true, nil
}

func (r repo) CreateOrderedProduct(ctx context.Context,
	orderID int64,
	orderRequests []model.OrderedProductDetail) error {
	if len(orderRequests) == 0 {
		return sql.ErrNoRows
	}
	query := `INSERT INTO ordered_products(
		product_id,
		order_id,
		qty,
		price_product,
		total_normal_price,
		total_final_price,
		discount_id)
		VALUES %s;`
	var values []interface{}
	for _, item := range orderRequests {
		values = append(values,
			item.ProductId,
			orderID,
			item.Qty,
			item.Price,
			item.TotalNormalPrice,
			item.TotalFinalPrice,
			item.DiscountID,
		)
	}
	template := "(?,?,?,?,?,?,?)"
	if len(orderRequests) > 1 {
		template += strings.Repeat(",(?,?,?,?,?,?,?)", len(orderRequests)-1)
	}
	query = fmt.Sprintf(query, template)
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(values...)
	if err != nil {
		return err
	}
	return nil
}
