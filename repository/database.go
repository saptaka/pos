package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DB interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type db struct {
	sqlDB *sql.DB
}

func newDatabase(username, password, host, dbName string,
	port, maxOpen, maxIdle int, timeout time.Duration) DB {
	log.Println("Connecting database ...")
	connection := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", username, password, host,
		port, dbName)

	dbConnection, err := sql.Open("mysql", connection)
	if err != nil {
		log.Fatal(err)
	}
	dbConnection.SetMaxIdleConns(maxIdle)
	dbConnection.SetMaxOpenConns(maxOpen)
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(time.Second*timeout))
	defer cancel()
	err = dbConnection.PingContext(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return db{dbConnection}
}

func (d db) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return d.sqlDB.QueryContext(ctx, query, args...)
}

func (d db) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return d.sqlDB.QueryRowContext(ctx, query, args...)
}

func (d db) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return d.sqlDB.ExecContext(ctx, query, args...)
}

func (d db) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return d.sqlDB.PrepareContext(ctx, query)
}

func (d db) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return d.sqlDB.BeginTx(ctx, opts)
}
