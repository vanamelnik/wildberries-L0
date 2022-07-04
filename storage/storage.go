package storage

// package storage describes the Storage interface and storage errors.

import (
	"errors"
)

type (
	// Storage represents app's order storage.
	Storage interface {
		Store(orderUID, jsonOrder string) error
		Get(orderUID string) (string, error)
		GetAll() ([]OrderDB, error)
	}

	// OrderDB represents the row in the database for order storing.
	OrderDB struct {
		OrderUID  string
		JSONOrder string
	}
)

var (
	ErrAlreadyExists = errors.New("order already exists")
	ErrNotFound      = errors.New("order not found")
)
