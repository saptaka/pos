package model

import "time"

type Product struct {
	ProductId  int64      `json:"productId"`
	Name       string     `json:"name" validate:"required"`
	SKU        string     `json:"sku,,omitempty"`
	Stock      int        `json:"stock,omitempty" validate:"required"`
	Price      int        `json:"price" validate:"required"`
	Image      string     `json:"image,omitempty"`
	CategoryID int        `json:"categoryId,omitempty"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
	CreatedAt  *time.Time `json:"createdAt,omitempty"`
	Discount   *Discount  `json:"discount,omitempty"`
	Category   *Category  `json:"category,omitempty"`
}

type Discount struct {
	DiscountID      int         `json:"discountId,omitempty"`
	Qty             int         `json:"qty" validate:"required"`
	Type            string      `json:"type" validate:"required"`
	Result          int         `json:"result"`
	ExpiratedAt     interface{} `json:"expiredAt,omitempty"`
	ExpiredAtFormat string      `json:"expiratedAtFormat,omitempty"`
	StringFormat    string      `json:"stringFormat,omitempty"`
}

type ListProduct struct {
	Products []Product `json:"products"`
	Meta     Meta      `json:"meta"`
}

var DiscountType = map[string]bool{
	"PERCENT": true,
	"BUY_N":   true,
}
