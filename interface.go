package main

import "todos/model"

// TodoService interface
type TodoService interface {
	All() ([]model.Todo, error)
	Create(todo *model.Todo) error
	FindByID(id int) (*model.Todo, error)
	DeleteByID(id int) ([]model.Todo, error)
	Update(id int, body string) (*model.Todo, error)
}

// SecretService interface
type SecretService interface {
	Insert(s *model.Secret) error
}
