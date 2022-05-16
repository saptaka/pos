package model

type Revenue struct {
	TotalRevenue int               `json:"totalRevenue"`
	PaymentType  []PaymentTypeItem `json:"paymentTypes"`
}

type PaymentTypeItem struct {
	Payment
	TotalAmount int `json:"totalAmount"`
}

type Solds struct {
	OrderProduct []SoldProduct `json:"orderProducts"`
}

type SoldProduct struct {
	ProductId   int64  `json:"productId"`
	Name        string `json:"name"`
	TotalQty    int    `json:"totalQty"`
	TotalAmount int    `json:"totalAmount"`
}
