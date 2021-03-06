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
	GetOrderByReceiptID(ctx context.Context, receiptId string) (model.Order, error)
	CreateOrder(ctx context.Context,
		orderRequest model.Order) (model.Order, error)
	DownloadReceipt(ctx context.Context, id int64) (string, error)
	GetDownloadStatus(ctx context.Context, id int64) (bool, error)
	CreateOrderedProduct(ctx context.Context, id int64, orderRequest []model.OrderedProductDetail) error
	GetOrderedProductByOrderId(ctx context.Context,
		id int64) ([]model.OrderedProductDetail, error)
}

func (r repo) GetOrder(ctx context.Context, limit, skip int) ([]model.Order, error) {

	ordersChan := make(chan []model.Order)

	go func(orderChanData chan []model.Order) {
		var rows *sql.Rows
		var err error
		query := `
			SELECT 
			id,
			payment_type_id,
			cashier_id,
			total_price,
			total_paid,
			total_return,
			receipt_id,
			created_at,
			updated_at 
			FROM 
			orders `
		if limit != 0 {
			query += " LIMIT ? OFFSET ?;"
			rows, err = r.db.QueryContext(ctx, query, limit, skip)
			if err != nil {
				orderChanData <- make([]model.Order, 0)
				return
			}
		} else {
			rows, err = r.db.QueryContext(ctx, query)
			if err != nil {
				orderChanData <- make([]model.Order, 0)
				return
			}
		}
		var orders []model.Order
		for rows.Next() {
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
				&order.UpdatedAt,
			)
			if err != nil {
				orderChanData <- make([]model.Order, 0)
				return
			}
			orders = append(orders, order)
		}
		orderChanData <- orders
	}(ordersChan)

	cashierChan := make(chan []model.Cashier)
	go func(cashierChanData chan []model.Cashier) {
		cashiers, err := r.GetCashiers(ctx, 0, 0)
		if err != nil {
			cashierChan <- make([]model.Cashier, 0)
			return
		}
		cashierChanData <- cashiers

	}(cashierChan)

	paymentChan := make(chan []model.Payment)
	go func(paymentChanData chan []model.Payment) {
		payments, err := r.GetPayments(ctx, 0, 0)
		if err != nil {
			paymentChanData <- make([]model.Payment, 0)
			return
		}
		paymentChanData <- payments
	}(paymentChan)

	orders := <-ordersChan
	cashiers := <-cashierChan
	mapCashier := make(map[int64]model.Cashier)
	for _, cashier := range cashiers {
		mapCashier[cashier.CashierId] = cashier
		cashiers = cashiers[1:]
	}

	payments := <-paymentChan
	mapPayment := make(map[int64]model.Payment)
	for _, payment := range payments {
		mapPayment[payment.PaymentId] = payment
		payments = payments[1:]
	}
	for index, order := range orders {
		if order.CashierID != nil {
			cashier := mapCashier[*order.CashierID]
			orders[index].Cashier = &cashier
		} else {
			cashierId := int64(0)
			orders[index].CashierID = &cashierId
		}
		if order.PaymentID != nil {
			payment := mapPayment[*order.PaymentID]
			orders[index].PaymentType = &payment
		}
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
		created_at,
		updated_at
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
		&order.UpdatedAt,
	)
	if err != nil {
		return order, err
	}

	cashierChan := make(chan model.Cashier)
	if order.CashierID != nil {
		go func(id int64, cashierChanData chan model.Cashier) {
			cashier, err := r.GetCashierByID(ctx, id)
			if err != nil {
				cashierChanData <- model.Cashier{}
				log.Println("error selecting cashier found")
				return
			}
			cashierChanData <- cashier

		}(*order.CashierID, cashierChan)
	} else {
		cashierId := int64(0)
		order.CashierID = &cashierId
		close(cashierChan)
	}

	paymentChan := make(chan model.Payment)
	if order.PaymentID != nil {
		go func(id int64, paymentChanData chan model.Payment) {
			payment, err := r.GetPaymentByID(ctx, id)
			if err != nil {
				paymentChanData <- model.Payment{}
				log.Println("error get payment found")
				return
			}
			paymentChanData <- payment
		}(*order.PaymentID, paymentChan)
	} else {
		close(paymentChan)
	}

	cashier := <-cashierChan
	order.Cashier = &cashier
	payment := <-paymentChan
	order.PaymentType = &payment

	return order, nil
}

func (r repo) GetOrderByReceiptID(ctx context.Context, receiptId string) (model.Order, error) {

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
		WHERE receipt_id=?;`

	rows := r.db.QueryRowContext(ctx, querySelect, receiptId)
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
	if err != nil {
		return order, err
	}

	cashierChan := make(chan model.Cashier)
	if order.CashierID != nil {
		go func(id int64, cashierChanData chan model.Cashier) {
			cashier, err := r.GetCashierByID(ctx, id)
			if err != nil {
				cashierChanData <- model.Cashier{}
				log.Println("error get cashier ", err)
				return
			}
			cashierChanData <- cashier

		}(*order.CashierID, cashierChan)
	} else {
		cashierId := int64(0)
		order.CashierID = &cashierId
		close(cashierChan)
	}

	paymentChan := make(chan model.Payment)
	if order.PaymentID != nil {
		go func(id int64, paymentChanData chan model.Payment) {
			payment, err := r.GetPaymentByID(ctx, id)
			if err != nil {
				paymentChanData <- model.Payment{}
				log.Println("error get payment ", err)
				return
			}
			paymentChanData <- payment
		}(*order.PaymentID, paymentChan)
	} else {
		close(cashierChan)
	}
	cashier := <-cashierChan
	order.Cashier = &cashier
	payment := <-paymentChan
	order.PaymentType = &payment

	return order, nil
}

func (r repo) CreateOrder(ctx context.Context, orderRequest model.Order) (model.Order, error) {

	query := `INSERT INTO orders(payment_type_id, total_price, total_paid, total_return, created_at, receipt_id)
			VALUES (?,?,?,?,?,?);`
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
		orderRequest.ReceiptID,
	)
	if err != nil {
		return orderRequest, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return orderRequest, err
	}
	orderRequest.OrderId = id
	if orderRequest.CashierID == nil {
		cashierId := int64(0)
		orderRequest.CashierID = &cashierId
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
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
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

func (r repo) CreateOrderedProduct(ctx context.Context, id int64, orderRequest []model.OrderedProductDetail) error {
	if len(orderRequest) == 0 {
		return sql.ErrNoRows
	}
	query := `INSERT INTO ordered_products(
		product_id,
		order_id,
		qty,
		price_product,
		name_product,
		total_normal_price,
		total_final_price,
		discount_id)
		VALUES %s;`
	var values []interface{}
	for _, item := range orderRequest {
		values = append(values,
			item.ProductId,
			id,
			item.Qty,
			item.Price,
			item.Name,
			item.TotalNormalPrice,
			item.TotalFinalPrice,
			item.DiscountId,
		)
	}
	template := "(?,?,?,?,?,?,?,?)"
	if len(orderRequest) > 1 {
		template += strings.Repeat(",(?,?,?,?,?,?,?,?)", len(orderRequest)-1)
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

func (r repo) GetOrderedProductByOrderId(ctx context.Context,
	id int64) ([]model.OrderedProductDetail, error) {

	query := `
	SELECT product_id,
		qty,
		total_normal_price,
		total_final_price,
		discount_id,
		price_product,
		name_product 
	FROM ordered_products
	WHERE order_id=?
	`
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	var orderedProducts []model.OrderedProductDetail
	for rows.Next() {
		var orderedProduct model.OrderedProductDetail
		err := rows.Scan(&orderedProduct.ProductId,
			&orderedProduct.Qty,
			&orderedProduct.TotalNormalPrice,
			&orderedProduct.TotalFinalPrice,
			&orderedProduct.DiscountId,
			&orderedProduct.Price,
			&orderedProduct.Name,
		)
		if err != nil {
			return nil, err
		}
		if orderedProduct.DiscountId != nil {
			discount, err := r.GetDiscountByID(ctx, *orderedProduct.DiscountId)
			if err != nil {
				return nil, err
			}
			orderedProduct.Discount = &discount
		}
		orderedProducts = append(orderedProducts, orderedProduct)
	}

	return orderedProducts, nil
}
