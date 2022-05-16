package model

import (
	"time"
)

type Cashier struct {
	ChashierId int64      `json:"cashierId,omitempty"`
	Name       string     `json:"name,omitempty"`
	Passcode   string     `json:"passcode,omitempty" validate:"required,len=6" `
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
	CreatedAt  *time.Time `json:"createdAt,omitempty"`
}

type ListCashier struct {
	Cashiers []Cashier `json:"cashiers"`
	Meta     Meta      `json:"meta"`
}
