package model

type Revenue struct {
	TotalRevenue int               `json:"totalRevenue"`
	PaymentType  []PaymentTypeItem `json:"paymentTypes"`
}

type PaymentTypeItem struct {
	PaymentTypeID int    `json:"paymentTypeId"`
	Name          int    `json:"name"`
	Logo          string `json:"logo"`
	TotalAmount   int    `json:"totalAmount"`
}
