package main

import "database/sql"

// SecretServiceImp is struct
type SecretServiceImp struct {
	db *sql.DB
}

//TodoServiceImp is struct
type TodoServiceImp struct {
	db *sql.DB
}
