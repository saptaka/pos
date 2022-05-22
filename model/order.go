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
	Name             string    `json:"name" `
	Price            int       `json:"price" `
	Discount         *Discount `json:"discount"`
	Qty              int       `json:"qty" `
	TotalNormalPrice int       `json:"totalNormalPrice"`
	TotalFinalPrice  int       `json:"totalFinalPrice"`
	DiscountId       *int64    `json:"-"`
}

type SubOrderedProductDetail struct {
	Product
	Qty              int `json:"qty" `
	TotalNormalPrice int `json:"totalNormalPrice"`
	TotalFinalPrice  int `json:"totalFinalPrice"`
}

type AddOrderRequest struct {
	PaymentID      int64            `json:"paymentId" `
	TotalPaid      int              `json:"totalPaid" `
	OrderedProduct []OrderedProduct `json:"products"`
}

type OrderedProduct struct {
	ProductId int64 `json:"productId" `
	Qty       int   `json:"qty" `
}

type OrderedProducts struct {
	Products []OrderedProduct
}

type SubTotalOrder struct {
	Subtotal       int                       `json:"subtotal"`
	OrderedProduct []SubOrderedProductDetail `json:"products"`
}
