package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var (
	db  *sql.DB
	err error
)

func main() {

	db, err = ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}

	createTable := CreateTable()
	if _, err := db.Exec(createTable); err != nil {
		log.Fatal(err)
	}

	services := &Server{
		db: db,
		todoService: &TodoServiceImp{
			db: db,
		},
		secretService: &SecretServiceImp{
			db: db,
		},
	}

	r := initializeRoutes(services)
	r.Run(":" + os.Getenv("PORT"))
}

//ConnectToDB is ConnectDatabase
func ConnectToDB() (*sql.DB, error) {

	database, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	return database, nil
}

//CreateTable is CreateTableDB
func CreateTable() string {
	return `
	CREATE TABLE IF NOT EXISTS todos (
		id SERIAL PRIMARY KEY,
		todo TEXT,
		created_at TIMESTAMP WITHOUT TIME ZONE,
		updated_at TIMESTAMP WITHOUT TIME ZONE
	);
	CREATE TABLE IF NOT EXISTS secrets (
		id SERIAL PRIMARY KEY,
		key TEXT
	);
	`
}
