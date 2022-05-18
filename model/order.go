package model

import "time"

type Orders struct {
	Order          Order                  `json:"order"`
	OrderedProduct []OrderedProductDetail `json:"products,omitempty"`
}

type Order struct {
	OrderId           int64      `json:"orderId"`
	PaymentID         *int64     `json:"paymentTypesId"`
	CashierID         *int64     `json:"cashierId"`
	TotalPaid         int        `json:"totalPaid"`
	TotalPrice        int        `json:"totalPrice"`
	TotalReturn       int        `json:"totalReturn"`
	ReceiptID         string     `json:"receiptId"`
	ReceiptIDFilePath string     `json:"-"`
	CreatedAt         *time.Time `json:"createdAt"`
	UpdatedAt         *time.Time `json:"updatedAt"`
	Cashier           Cashier    `json:"cashier"`
	PaymentType       Payment    `json:"payment_type"`
}

type OrderedProductDetail struct {
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
	Subtotal       int                    `json:"subtotal"`
	OrderedProduct []OrderedProductDetail `json:"products"`
}
