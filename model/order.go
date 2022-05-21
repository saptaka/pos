package model

import "time"

type OrderDetails struct {
	Order          Order                  `json:"order"`
	OrderedProduct []OrderedProductDetail `json:"products,omitempty"`
}

type ListOrders struct {
	Order []Order `json:"orders"`
	Meta  Meta    `json:"meta"`
}

type Order struct {
	OrderId           int64      `json:"orderId"`
	CashierID         *int64     `json:"cashiersId,omitempty"`
	PaymentID         *int64     `json:"paymentTypesId"`
	TotalPrice        int        `json:"totalPrice"`
	TotalPaid         int        `json:"totalPaid"`
	TotalReturn       int        `json:"totalReturn"`
	ReceiptID         string     `json:"receiptId"`
	ReceiptIDFilePath string     `json:"-"`
	UpdatedAt         *time.Time `json:"updatedAt"`
	CreatedAt         *time.Time `json:"createdAt"`
	Cashier           *Cashier   `json:"cashier,omitempty"`
	PaymentType       *Payment   `json:"payment_type,omitempty"`
}

type OrderedProductDetail struct {
	ProductId        int64     `json:"productId"`
	Name             string    `json:"name" validate:"required"`
	Price            int       `json:"price" validate:"required"`
	Discount         *Discount `json:"discount"`
	Qty              int       `json:"qty" validate:"required"`
	TotalFinalPrice  int       `json:"totalFinalPrice"`
	TotalNormalPrice int       `json:"totalNormalPrice"`
	DiscountId       *int64    `json:"-"`
}

type SubOrderedProductDetail struct {
	Product
	Qty              int `json:"qty" validate:"required"`
	TotalFinalPrice  int `json:"totalFinalPrice"`
	TotalNormalPrice int `json:"totalNormalPrice"`
}

type AddOrderRequest struct {
	PaymentID      int64            `json:"paymentId" validate:"required"`
	TotalPaid      int              `json:"totalPaid" validate:"required"`
	OrderedProduct []OrderedProduct `json:"products"`
}

type OrderedProduct struct {
	ProductId int64 `json:"productId" validate:"required"`
	Qty       int   `json:"qty" validate:"required"`
}

type SubTotalOrder struct {
	Subtotal       int                       `json:"subtotal"`
	OrderedProduct []SubOrderedProductDetail `json:"products"`
}
