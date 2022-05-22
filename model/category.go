package model

import "time"

type Category struct {
	CategoryId int64      `json:"categoryId"`
	Name       string     `json:"name"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
	CreatedAt  *time.Time `json:"createdAt,omitempty"`
}

type ListCategory struct {
	Categories []Category `json:"categories"`
	Meta       Meta       `json:"meta"`
}
