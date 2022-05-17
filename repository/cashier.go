package repository

import (
	"context"
	"database/sql"

	"github.com/saptaka/pos/model"
)

type CashierRepo interface {
	GetCashierByID(ctx context.Context, id int64) (model.Cashier, error)
	GetCashiers(ctx context.Context, limit, skip int) ([]model.Cashier, error)
	UpdateCashier(ctx context.Context, cashier model.Cashier) error
	CreateCashier(ctx context.Context, name, passcode string) (model.Cashier, error)
	DeleteCashier(ctx context.Context, id int64) error
	GetPasscodeById(ctx context.Context, id int64) (string, error)
}

func (r repo) GetCashierByID(ctx context.Context, id int64) (model.Cashier, error) {
	var cashier model.Cashier
	query := "SELECT id, name FROM cashiers WHERE id=?"
	rows := r.db.QueryRowContext(ctx, query, id)
	err := rows.Scan(&cashier.CashierId, &cashier.Name)
	if err != nil {
		return cashier, err
	}

	return cashier, nil
}

func (r repo) GetCashiers(ctx context.Context,
	limit, skip int) ([]model.Cashier, error) {
	query := "SELECT id, name FROM cashiers "
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

	var cashiers []model.Cashier
	for rows.Next() {
		var cashier model.Cashier
		err := rows.Scan(&cashier.CashierId, &cashier.Name)
		if err != nil {
			return nil, err
		}
		cashiers = append(cashiers, cashier)
	}
	return cashiers, nil
}

func (r repo) UpdateCashier(ctx context.Context,
	cashierDetail model.Cashier) error {
	query := `UPDATE cashiers 
		SET name=?, 
			passcode=?, 
			updated_at=CURRENT_TIMESTAMP() 
		WHERE id=?`
	result, err := r.db.ExecContext(ctx, query,
		cashierDetail.Name,
		cashierDetail.Passcode,
		cashierDetail.CashierId)
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

func (r repo) CreateCashier(ctx context.Context, name, passcode string) (model.Cashier, error) {
	var cashier model.Cashier
	insertQuery := `INSERT INTO 
		cashiers (name,passcode) 
	VALUES (?,?);`
	stmt, err := r.db.PrepareContext(ctx, insertQuery)
	if err != nil {
		return cashier, err
	}
	res, err := stmt.Exec(name, passcode)
	if err != nil {
		return cashier, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return cashier, err
	}
	selectQuery := `SELECT id, 
					name,
					passcode, 
					updated_at, 
					created_at
					FROM cashiers 
					WHERE id=?`
	rows := r.db.QueryRowContext(ctx, selectQuery, id)

	err = rows.Scan(
		&cashier.CashierId,
		&cashier.Name,
		&cashier.Passcode,
		&cashier.UpdatedAt,
		&cashier.CreatedAt)
	return cashier, err

}

func (r repo) DeleteCashier(ctx context.Context, id int64) error {
	_, err := r.GetCashierByID(ctx, id)
	if err == sql.ErrNoRows {
		return sql.ErrNoRows
	}

	if err != nil {
		return err
	}
	query := "DELETE FROM cashiers WHERE id=?"
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

func (r repo) GetPasscodeById(ctx context.Context, id int64) (string, error) {
	var passcode string
	query := "SELECT passcode FROM cashiers WHERE id=?"
	rows := r.db.QueryRowContext(ctx, query, id)
	err := rows.Scan(&passcode)
	if err != nil {
		return passcode, err
	}

	return passcode, nil
}
