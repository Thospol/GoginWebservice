package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Server struct {
	db *sql.DB
}

type Todo struct {
	ID        int64     `json:"id"`
	Body      string    `json:"todo"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (s *Server) FindByID(c *gin.Context) {
	name := c.Param("id")
	var todo Todo
	queryStmt := "SELECT id, todo, updated_at, created_at FROM todos where id = $1"
	row := s.db.QueryRow(queryStmt, name)

	err := row.Scan(&todo.ID, &todo.Body, &todo.UpdatedAt, &todo.CreatedAt)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, todo)

}

func (s *Server) All(c *gin.Context) {
	rows, err := s.db.Query("SELECT id, todo, updated_at, created_at FROM todos")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"object":  "error",
			"message": fmt.Sprintf("db: query error: %s", err),
		})
		return
	}
	todos := []Todo{} // set empty slice without nil
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Body, &todo.UpdatedAt, &todo.CreatedAt)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"object":  "error",
				"message": fmt.Sprintf("db: query error: %s", err),
			})
			return
		}
		todos = append(todos, todo)
	}
	c.JSON(http.StatusOK, todos)
}

func (s *Server) Create(c *gin.Context) {
	var todo Todo
	err := c.ShouldBindJSON(&todo)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"object":  "error",
			"message": fmt.Sprintf("json: wrong params: %s", err),
		})
		return
	}
	now := time.Now()
	todo.CreatedAt = now
	todo.UpdatedAt = now
	row := s.db.QueryRow("INSERT INTO todos (todo, created_at, updated_at) values ($1, $2, $3) RETURNING id", todo.Body, now, now)

	if err := row.Scan(&todo.ID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"object":  "error",
			"message": fmt.Sprintf("db: query error: %s", err),
		})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	createTable := `
	CREATE TABLE IF NOT EXISTS todos (
		id SERIAL PRIMARY KEY,
		todo TEXT,
		created_at TIMESTAMP WITHOUT TIME ZONE,
		updated_at TIMESTAMP WITHOUT TIME ZONE
	);
	`
	if _, err := db.Exec(createTable); err != nil {
		log.Fatal(err)
	}

	s := &Server{
		db: db,
	}
	r := gin.Default()
	r.GET("/todos", s.All)
	r.POST("/todos", s.Create)
	r.GET("/todo/:name", s.FindByID)

	r.Run(":" + os.Getenv("PORT"))
}
