package main

import (
	"time"
	"todos/model"
)

//All is Medthod of TodoServiceImp
func (s *TodoServiceImp) All() ([]model.Todo, error) {
	rows, err := s.db.Query("SELECT id, todo, updated_at, created_at FROM todos")
	if err != nil {
		return nil, err
	}
	todos := []model.Todo{} // set empty slice without nil
	for rows.Next() {
		var todo model.Todo
		err := rows.Scan(&todo.ID, &todo.Body, &todo.UpdatedAt, &todo.CreatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

//Create is Medthod of TodoServiceImp
func (s *TodoServiceImp) Create(todo *model.Todo) error {
	now := time.Now()
	todo.CreatedAt = now
	todo.UpdatedAt = now
	row := s.db.QueryRow("INSERT INTO todos (todo, created_at, updated_at) values ($1, $2, $3) RETURNING id", todo.Body, now, now)

	if err := row.Scan(&todo.ID); err != nil {
		return err
	}
	return nil
}

//FindByID is Medthod of TodoServiceImp
func (s *TodoServiceImp) FindByID(id int) (*model.Todo, error) {
	stmt := "SELECT id, todo, created_at, updated_at FROM todos WHERE id = $1"
	row := s.db.QueryRow(stmt, id)
	var todo model.Todo
	err := row.Scan(&todo.ID, &todo.Body, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

//DeleteByID is Medthod of TodoServiceImp
func (s *TodoServiceImp) DeleteByID(id int) ([]model.Todo, error) {

	stmt := "DELETE FROM todos WHERE id = $1"
	_, err := s.db.Exec(stmt, id)
	if err != nil {
		return nil, err
	}
	todos, err := s.All()
	if err != nil {
		return nil, err
	}
	return todos, nil
}

//Update is Medthod of TodoServiceImp
func (s *TodoServiceImp) Update(id int, body string) (*model.Todo, error) {
	stmt := "UPDATE todos SET todo = $2 WHERE id = $1"
	_, err := s.db.Exec(stmt, id, body)
	if err != nil {
		return nil, err
	}
	return s.FindByID(id)
}
