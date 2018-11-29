package model

import "time"

//Todo is Model Todo
type Todo struct {
	ID        int64     `json:"id"`
	Body      string    `json:"todo" binding:"required"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

//Secret is Model Secret
type Secret struct {
	ID  int64  `json:"id"`
	Key string `json:"key" binding:"required"`
}
