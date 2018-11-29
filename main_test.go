package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"todos/model"

	"github.com/stretchr/testify/assert"
)

type fakeTodoService struct {
}

func (s *fakeTodoService) All() ([]model.Todo, error) {
	return []model.Todo{
		{
			ID:        1,
			Body:      "Todo",
			CreatedAt: time.Date(2018, 11, 29, 13, 26, 0, 0, time.Local),
			UpdatedAt: time.Date(2018, 11, 29, 13, 26, 0, 0, time.Local),
		},
	}, nil
}

func TestAll(t *testing.T) {
	services := &Server{
		//todoService: &fakeTodoService{},
	}
	r := initializeRoutes(services)
	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/todos", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusOK)
}
