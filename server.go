package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"todos/model"

	"github.com/gin-gonic/gin"
)

// Server struct
type Server struct {
	db            *sql.DB
	todoService   TodoService
	secretService SecretService
}

//FindByID is Medthod of Server
func (s *Server) FindByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	todo, err := s.todoService.FindByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, todo)
}

//All is Medthod of Server
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

//Create is Medthod of Server
func (s *Server) Create(c *gin.Context) {

	var todo model.Todo
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

//Update is Medthod of Server
func (s *Server) Update(c *gin.Context) {

	h := map[string]string{}
	err := c.ShouldBindJSON(&h)
	if err != nil {
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

//DeleteByID is Medthod of Server
func (s *Server) DeleteByID(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))
	todos, err := s.todoService.DeleteByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, todos)

}

//CreateSecret is Medthod of Server
func (s *Server) CreateSecret(c *gin.Context) {
	var secret model.Secret
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

//AuthTodo is Medthod of Server
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
