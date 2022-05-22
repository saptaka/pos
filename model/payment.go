package model

import "time"

type Payment struct {
	PaymentId int64      `json:"paymentId"`
	Name      string     `json:"name" `
	Type      string     `json:"type" `
	Logo      string     `json:"logo"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
}

type ListPayment struct {
	Payments []Payment `json:"payments"`
	Meta     Meta      `json:"meta"`
}

var PaymentType = map[string]bool{
	"CASH":     true,
	"E-WALLET": true,
	"EDC":      true,
}
