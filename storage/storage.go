package storage

import (
	"errors"
)

type Storage interface {
	Store(orderUID, jsonOrder string) error
	Get(orderUID string) (string, error)
	GetAll() ([]string, error)
}

var (
	ErrAlreadyExists = errors.New("order already exists")
	ErrNotFound      = errors.New("order not found")
)
