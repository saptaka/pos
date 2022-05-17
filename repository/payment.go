package repository

import (
	"context"
	"database/sql"

	"github.com/saptaka/pos/model"
)

type PaymentRepo interface {
	GetPaymentByID(ctx context.Context, id int64) (model.Payment, error)
	GetPayments(ctx context.Context, limit, skip int) ([]model.Payment, error)
	UpdatePayment(ctx context.Context, payment model.Payment) error
	CreatePayment(ctx context.Context, payment model.Payment) (model.Payment, error)
	DeletePayment(ctx context.Context, id int) error
}

func (r repo) GetPaymentByID(ctx context.Context, id int64) (model.Payment, error) {
	var payment model.Payment
	query := "SELECT id, name, types,logo FROM payments WHERE id=?"
	rows := r.db.QueryRowContext(ctx, query, id)
	err := rows.Scan(&payment.PaymentId, &payment.Name,
		&payment.Type, &payment.Logo)

	if err == sql.ErrNoRows {
		return payment, err
	}
	if err != nil {
		return payment, err
	}

	return payment, nil
}

func (r repo) GetPayments(ctx context.Context,
	limit, skip int) ([]model.Payment, error) {
	query := "SELECT id, name, types, logo FROM payments "
	var rows *sql.Rows
	var err error
	if limit > 0 {
		query += " limit ? offset ?;"
		rows, err = r.db.QueryContext(ctx, query, limit, skip)
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		if err != nil {
			return nil, err
		}
	} else {
		rows, err = r.db.QueryContext(ctx, query)
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		if err != nil {
			return nil, err
		}
	}
	var payments []model.Payment
	for rows.Next() {
		var payment model.Payment
		err := rows.Scan(&payment.PaymentId,
			&payment.Name,
			&payment.Type,
			&payment.Logo)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	return payments, nil
}

func (r repo) UpdatePayment(ctx context.Context,
	payment model.Payment) error {
	query := "UPDATE payments SET name=?, types=?, logo=? ,updated_at=CURRENT_TIMESTAMP() WHERE id=?"
	result, err := r.db.ExecContext(ctx, query, payment.Name, payment.Type, payment.Logo,
		payment.PaymentId)
	if err != nil {
		return sql.ErrNoRows
	}
	rowaffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowaffected == 0 {
		return sql.ErrNoRows
	}
	return err
}

func (r repo) CreatePayment(ctx context.Context, payment model.Payment) (model.Payment, error) {
	var paymentRequest model.Payment
	insertQuery := `INSERT INTO 
		payments (name, types, logo) 
	VALUES (?,?,?);`
	stmt, err := r.db.PrepareContext(ctx, insertQuery)
	if err != nil {
		return paymentRequest, err
	}
	res, err := stmt.Exec(payment.Name, payment.Type, payment.Logo)
	if err != nil {
		return paymentRequest, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return paymentRequest, err
	}
	selectQuery := `SELECT id, 
					name,
					types,
					logo, 
					updated_at, 
					created_at
					FROM payments 
					WHERE id=?`
	rows := r.db.QueryRowContext(ctx, selectQuery, id)

	err = rows.Scan(
		&paymentRequest.PaymentId,
		&paymentRequest.Name,
		&paymentRequest.Type,
		&paymentRequest.Logo,
		&paymentRequest.UpdatedAt,
		&paymentRequest.CreatedAt)
	return paymentRequest, err

}

func (r repo) DeletePayment(ctx context.Context, id int) error {
	query := "DELETE FROM payments WHERE id=?"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
