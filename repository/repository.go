package repository

import (
	"time"

	"github.com/saptaka/pos/config"
)

type Repo interface {
	CashierRepo
	CategoryRepo
	ProductRepo
	PaymentRepo
	OrderRepo
	ReportRepo
	RevenueRepo
}

type repo struct {
	db DB
}

func NewRepository(cfg *config.Config) Repo {
	database := newDatabase(cfg.DBuser, cfg.DBPassword, cfg.DBHost, cfg.DBName,
		cfg.DBPort, cfg.DBMaxConnection, cfg.DBMaxIdle,
		time.Duration(cfg.DBConnectionTimeout))
	return &repo{database}
}
