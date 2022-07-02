package inmem

import (
	"sync"

	"github.com/vanamelnik/wildberries-L0/storage"
)

var _ storage.Storage = (*Storage)(nil)

type Storage struct {
	mu         *sync.RWMutex
	repository map[string]string
}

func NewStorage() *Storage {
	return &Storage{
		mu:         &sync.RWMutex{},
		repository: make(map[string]string),
	}
}

func (s *Storage) Store(orderUID, jsonOrder string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.repository[orderUID]; ok {
		return storage.ErrAlreadyExists
	}
	s.repository[orderUID] = jsonOrder
	return nil
}

func (s *Storage) Get(orderUID string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	order, ok := s.repository[orderUID]
	if !ok {
		return "", storage.ErrNotFound
	}
	return order, nil
}

func (s *Storage) GetAll() ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	orders := make([]string, 0, len(s.repository))
	for _, order := range s.repository {
		orders = append(orders, order)
	}
	return orders, nil
}
