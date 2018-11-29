package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type TodoService interface {
	All() ([]Todo, error)
	Create(todo *Todo) error
	FindByID(id int) (*Todo, error)
	DeleteByID(id int) error
	Update(id int, body string) (*Todo, error)
}

// type Server struct {
// 	db      *sql.DB
// 	service TodoService
// }

type Server struct {
	db            *sql.DB
	todoService   TodoService
	secretService SecretService
}

func (s *Server) FindByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	todo, err := s.todoService.FindByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, todo)
}

func (s *Server) All(c *gin.Context) {
	todos, err := s.todoService.All()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"object":  "error",
			"message": fmt.Sprintf("db: query error: %s", err),
		})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func (s *Server) Create(c *gin.Context) {
	// var todo Todo
	// err := c.ShouldBindJSON(&todo)
	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
	// 		"object":  "error",
	// 		"message": fmt.Sprintf("json: wrong params: %s", err),
	// 	})
	// 	return
	// }
	// now := time.Now()
	// todo.CreatedAt = now
	// todo.UpdatedAt = now
	// row := s.db.QueryRow("INSERT INTO todos (todo, created_at, updated_at) values ($1, $2, $3) RETURNING id", todo.Body, now, now)

	// if err := row.Scan(&todo.ID); err != nil {
	// 	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
	// 		"object":  "error",
	// 		"message": fmt.Sprintf("db: query error: %s", err),
	// 	})
	// 	return
	// }

	// c.JSON(http.StatusCreated, todo)

	var todo Todo
	err := c.ShouldBindJSON(&todo)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"object":  "error",
			"message": fmt.Sprintf("json: wrong params: %s", err),
		})
		return
	}

	if err := s.todoService.Create(&todo); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, todo)
}

func (s *Server) Update(c *gin.Context) {

	// h := gin.H{}
	// if err := c.ShouldBindJSON(&h); err != nil {
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, err)
	// 	return
	// }
	// id, _ := strconv.Atoi(c.Param("id"))
	// stmt := "UPDATE todos SET todo = $2 WHERE id = $1"
	// _, err := s.db.Exec(stmt, id, h["todo"])
	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusInternalServerError, err)
	// 	return
	// }
	// stmt = "SELECT id, todo, created_at, updated_at FROM todos WHERE id = $1"
	// row := s.db.QueryRow(stmt, id)
	// var todo Todo
	// err = row.Scan(&todo.ID, &todo.Body, &todo.CreatedAt, &todo.UpdatedAt)
	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusInternalServerError, err)
	// 	return
	// }
	// c.JSON(http.StatusOK, todo)

	h := map[string]string{}
	if err := c.ShouldBindJSON(&h); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	todo, err := s.todoService.Update(id, h["todo"])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, todo)
}

func (s *Server) DeleteByID(c *gin.Context) {
	// stmt := "DELETE FROM todos WHERE id = $1"
	// id, _ := strconv.Atoi(c.Param("id"))
	// _, err := s.db.Exec(stmt, id)
	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusInternalServerError, err)
	// 	return
	// }

	id, _ := strconv.Atoi(c.Param("id"))
	if err := s.todoService.DeleteByID(id); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
}

type TodoServiceImp struct {
	db *sql.DB
}

func (s *TodoServiceImp) All() ([]Todo, error) {
	rows, err := s.db.Query("SELECT id, todo, updated_at, created_at FROM todos")
	if err != nil {
		return nil, err
	}
	todos := []Todo{} // set empty slice without nil
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Body, &todo.UpdatedAt, &todo.CreatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func (s *TodoServiceImp) Create(todo *Todo) error {
	now := time.Now()
	todo.CreatedAt = now
	todo.UpdatedAt = now
	row := s.db.QueryRow("INSERT INTO todos (todo, created_at, updated_at) values ($1, $2, $3) RETURNING id", todo.Body, now, now)

	if err := row.Scan(&todo.ID); err != nil {
		return err
	}
	return nil
}

func (s *TodoServiceImp) FindByID(id int) (*Todo, error) {
	stmt := "SELECT id, todo, created_at, updated_at FROM todos WHERE id = $1"
	row := s.db.QueryRow(stmt, id)
	var todo Todo
	err := row.Scan(&todo.ID, &todo.Body, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (s *TodoServiceImp) DeleteByID(id int) error {
	stmt := "DELETE FROM todos WHERE id = $1"
	_, err := s.db.Exec(stmt, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *TodoServiceImp) Update(id int, body string) (*Todo, error) {
	stmt := "UPDATE todos SET todo = $2 WHERE id = $1"
	_, err := s.db.Exec(stmt, id, body)
	if err != nil {
		return nil, err
	}
	return s.FindByID(id)
}

type Todo struct {
	ID        int64     `json:"id"`
	Body      string    `json:"todo" binding:"required"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Secret struct {
	ID  int64  `json:"id"`
	Key string `json:"key" binding:"required"`
}

type SecretService interface {
	Insert(s *Secret) error
}
type SecretServiceImp struct {
	db *sql.DB
}

func (s *SecretServiceImp) Insert(secret *Secret) error {
	row := s.db.QueryRow("INSERT INTO secrets (key) values ($1) RETURNING id", secret.Key)
	if err := row.Scan(&secret.ID); err != nil {
		return err
	}
	return nil
}

func (s *Server) CreateSecret(c *gin.Context) {
	var secret Secret
	if err := c.ShouldBindJSON(&secret); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	if err := s.secretService.Insert(&secret); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, secret)
}

func (s *Server) AuthTodo(c *gin.Context) {

	user, _, ok := c.Request.BasicAuth()
	if ok {
		row := s.db.QueryRow("SELECT key FROM secrets WHERE key = $1", user)
		if err := row.Scan(&user); err == nil {
			return
		}
	}
	c.AbortWithStatus(http.StatusUnauthorized)
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
	CREATE TABLE IF NOT EXISTS secrets (
		id SERIAL PRIMARY KEY,
		key TEXT
	);
	`

	if _, err := db.Exec(createTable); err != nil {
		log.Fatal(err)
	}

	// s := &Server{
	// 	db: db,
	// }

	// s := &Server{
	// 	service: &TodoServiceImp{
	// 		db: db,
	// 	},
	// }

	s := &Server{
		db: db,
		todoService: &TodoServiceImp{
			db: db,
		},
		secretService: &SecretServiceImp{
			db: db,
		},
	}

	r := SetupRoute(s)
	r.Run(":" + os.Getenv("PORT"))
}

func SetupRoute(s *Server) *gin.Engine {

	r := gin.Default()
	todos := r.Group("/todos")
	admin := r.Group("/admin")
	admin.Use(gin.BasicAuth(gin.Accounts{
		"admin": "1234",
	}))

	// r.Use(func(c *gin.Context) {
	// 	user, pass, ok := c.Request.BasicAuth()
	// 	if ok {
	// 		if user == "foo" && pass == "pass" {
	// 			c.Set(gin.AuthUserKey, user)
	// 			return
	// 		}
	// 	}
	// 	c.AbortWithStatus(http.StatusUnauthorized)
	// })
	todos.Use(s.AuthTodo)
	todos.GET("/", s.All)
	todos.POST("/", s.Create)
	todos.GET("/:id", s.FindByID)
	todos.PUT("/:id", s.Update)
	todos.DELETE("/:id", s.DeleteByID)

	// curl -XPOST http://localhost:8000/admin/secrets
	//  -u admin:1234 -d '{"key": "foobar"}
	admin.POST("/admin/secrets", s.CreateSecret)
	return r
}
