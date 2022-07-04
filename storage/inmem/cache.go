package inmem

import (
	"fmt"
	"log"
	"sync"

	"github.com/vanamelnik/wildberries-L0/storage"
)

var _ storage.Storage = (*Cache)(nil)

type (
	Cache struct {
		mu                *sync.RWMutex
		repository        map[string]string
		persistentStorage storage.Storage
	}

	StorageOpt func(s *Cache) error
)

func NewCache(opts ...StorageOpt) (*Cache, error) {
	s := &Cache{
		mu:         &sync.RWMutex{},
		repository: make(map[string]string),
	}
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, fmt.Errorf("storage: inmem: could not apply option: %w", err)
		}
	}

	return s, nil
}

func WithPersistentStorage(ps storage.Storage) StorageOpt {
	return func(s *Cache) error {
		orders, err := ps.GetAll()
		if err != nil {
			return err
		}
		for _, o := range orders {
			s.Store(o.OrderUID, o.JSONOrder)
		}
		s.persistentStorage = ps
		if len(orders) > 0 {
			log.Printf("storage: inmem: %d record(s) successfully imported from the database", len(orders))
		}
		return nil
	}
}

func (s *Cache) Store(orderUID, jsonOrder string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.repository[orderUID]; ok {
		return storage.ErrAlreadyExists
	}
	s.repository[orderUID] = jsonOrder
	if s.persistentStorage != nil {
		s.persistentStorage.Store(orderUID, jsonOrder)
	}
	return nil
}

func (s *Cache) Get(orderUID string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	order, ok := s.repository[orderUID]
	if !ok {
		return "", storage.ErrNotFound
	}
	return order, nil
}

func (s *Cache) GetAll() ([]storage.OrderDB, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	orders := make([]storage.OrderDB, 0, len(s.repository))
	for uid, order := range s.repository {
		orders = append(orders, storage.OrderDB{
			OrderUID:  uid,
			JSONOrder: order,
		})
	}
	return orders, nil
}
