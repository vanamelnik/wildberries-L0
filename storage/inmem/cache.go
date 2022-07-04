package inmem

// inmem is in-memory cache that could be also used as independed repository.

import (
	"fmt"
	"log"
	"sync"

	"github.com/vanamelnik/wildberries-L0/storage"
)

var _ storage.Storage = (*Cache)(nil)

type (
	// Cache is an in-memory implementation of storage.Storage.
	// It could be used as indepened repository or use another storage.Storage object
	// for persistent storage.
	Cache struct {
		mu                *sync.RWMutex
		repository        map[string]string
		persistentStorage storage.Storage
	}

	StorageOpt func(s *Cache) error
)

// NewCache creates a new in-memory repository and registers persistent storage if provided.
func NewCache(opts ...StorageOpt) (*Cache, error) {
	s := &Cache{
		mu:         &sync.RWMutex{},
		repository: make(map[string]string),
	}
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, fmt.Errorf("storage: cache: could not apply option: %w", err)
		}
	}

	return s, nil
}

// WithPersistentStorage registers a given storage.Storage object as persistent storage.
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

// Store implements storage.Storage interface.
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

// Get implements storage.Storage interface.
func (s *Cache) Get(orderUID string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	order, ok := s.repository[orderUID]
	if !ok {
		return "", storage.ErrNotFound
	}
	return order, nil
}

// GetAll implements storage.Storage interface.
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
