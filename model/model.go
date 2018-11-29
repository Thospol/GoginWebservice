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

//User is Model User
type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	FirstName string `json:"firstname" binding:"required"`
	LastName  string `json:"lastname" binding:"required"`
	Age       int64  `json:"age" binding:"required"`
	IDCard    string `json:"idcard" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Telephone string `json:"telephone" binding:"required"`
}
