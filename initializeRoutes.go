package main

import "github.com/gin-gonic/gin"

func initializeRoutes(s *Server) *gin.Engine {

	r := gin.Default()
	todos := r.Group("/todos")
	admin := r.Group("/admin")
	admin.Use(gin.BasicAuth(gin.Accounts{
		"admin": "1234",
	}))

	todos.Use(s.AuthTodo)
	todos.GET("/", s.All)
	todos.POST("/", s.Create)
	todos.GET("/:id", s.FindByID)
	todos.PUT("/:id", s.Update)
	todos.DELETE("/:id", s.DeleteByID)

	admin.POST("/secrets", s.CreateSecret)

	return r
}
