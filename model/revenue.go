package model

type Revenue struct {
	TotalRevenue int               `json:"totalRevenue"`
	PaymentType  []PaymentTypeItem `json:"paymentTypes"`
}

type PaymentTypeItem struct {
	Payment
	TotalAmount int `json:"totalAmount"`
}
