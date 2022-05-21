package model

import "time"

type ProductCreateRequest struct {
	Name       string    `json:"name" validate:"required"`
	Stock      int       `json:"stock,omitempty" validate:"required"`
	Price      int       `json:"price" validate:"required"`
	Image      string    `json:"image,omitempty"`
	CategoryId *int64    `json:"categoryId"`
	Discount   *Discount `json:"discount"`
}

type Product struct {
	ProductId  int64      `json:"productId"`
	Name       string     `json:"name" validate:"required"`
	Stock      int        `json:"stock,omitempty" validate:"required"`
	Price      int        `json:"price" validate:"required"`
	Image      string     `json:"image,omitempty"`
	SKU        string     `json:"sku,omitempty"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
	CreatedAt  *time.Time `json:"createdAt,omitempty"`
	DiscountId *int64     `json:"discountId,omitempty"`
	CategoryId *int64     `json:"categoryId,omitempty"`
	Discount   *Discount  `json:"discount"`
	Category   *Category  `json:"category,omitempty"`
}

type ProductCreateResponse struct {
	ProductId  int64      `json:"productId"`
	Name       string     `json:"name" validate:"required"`
	Stock      int        `json:"stock,omitempty" validate:"required"`
	Price      int        `json:"price" validate:"required"`
	Image      string     `json:"image,omitempty"`
	SKU        string     `json:"sku,omitempty"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
	CreatedAt  *time.Time `json:"createdAt,omitempty"`
	CategoryId *int64     `json:"categoryId"`
}

type Discount struct {
	DiscountID      int64       `json:"discountId,omitempty"`
	Qty             int         `json:"qty" validate:"required"`
	Type            string      `json:"type" validate:"required"`
	Result          int         `json:"result"`
	ExpiratedAt     interface{} `json:"expiredAt,omitempty"`
	ExpiredAtFormat string      `json:"expiratedAtFormat"`
	StringFormat    string      `json:"stringFormat"`
}

type ListProduct struct {
	Products []Product `json:"products"`
	Meta     Meta      `json:"meta"`
}

var DiscountType = map[string]bool{
	"PERCENT": true,
	"BUY_N":   true,
}

const Percent = "PERCENT"
const BuyN = "BUY_N"
