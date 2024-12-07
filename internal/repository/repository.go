package repository

import (
	"playground/internal/entities"
)

type Config struct {
}

func NewRepository(config *Config) (*entities.Repository, error) {
	repo := &entities.Repository{}
	return repo, nil
}
