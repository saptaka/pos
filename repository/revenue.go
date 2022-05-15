package repository

import (
	"context"

	"github.com/saptaka/pos/model"
)

type RevenueRepo interface {
	GetRevenues(ctx context.Context) ([]model.Revenue, error)
}

func (r repo) GetRevenues(ctx context.Context) ([]model.Revenue, error) {

	return nil, nil
}
