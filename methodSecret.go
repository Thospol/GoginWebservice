package main

import (
	"todos/model"
)

//Insert is Medthod of SecretServiceImp
func (s *SecretServiceImp) Insert(secret *model.Secret) error {
	row := s.db.QueryRow("INSERT INTO secrets (key) values ($1) RETURNING id", secret.Key)
	if err := row.Scan(&secret.ID); err != nil {
		return err
	}
	return nil
}
