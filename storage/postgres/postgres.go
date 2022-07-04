package postgres

import (
	"database/sql"
	_ "embed"
	"errors"
	"log"
	"sync"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/vanamelnik/wildberries-L0/storage"
)

const storeChSize = 1000 // the size of storeCh buffer

type (
	// Storage is an implementation of storage.Storage using Postgresql db engine.
	// Saving the orders works in async mode.
	Storage struct {
		db      *sql.DB
		storeCh chan storage.OrderDB
		stopCh  chan struct{}
		wg      *sync.WaitGroup
	}
)

var _ storage.Storage = (*Storage)(nil)

//go:embed schema.sql
var queryCreate string

// NewStorage connects to the postgres database and runs the worker that stores all incoming orders.
func NewStorage(databaseURI string) (*Storage, error) {
	db, err := sql.Open("pgx", databaseURI)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if _, err := db.Exec(queryCreate); err != nil {
		return nil, err
	}
	s := &Storage{
		db:      db,
		storeCh: make(chan storage.OrderDB, storeChSize),
		stopCh:  make(chan struct{}),
		wg:      &sync.WaitGroup{},
	}
	s.wg.Add(1)
	go s.storer()
	return s, nil
}

// Close stops the worker and closes the db connection.
func (s *Storage) Close() error {
	if s.stopCh != nil {
		close(s.stopCh)
		s.stopCh = nil
	}
	s.wg.Wait()
	return s.db.Close()
}

// Get implements storage.Storage interface.
func (s *Storage) Get(orderUID string) (string, error) {
	var order string
	err := s.db.QueryRow(`SELECT json_order FROM orders WHERE uid = $1;`, orderUID).Scan(&order)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrNotFound
		}
		return "", err
	}
	return order, nil
}

// GetAll implements storage.Storage interface.
func (s *Storage) GetAll() ([]storage.OrderDB, error) {
	orders := make([]storage.OrderDB, 0)
	rows, err := s.db.Query("SELECT uid, json_order FROM orders;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var o storage.OrderDB
		if err := rows.Scan(&o.OrderUID, &o.JSONOrder); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

// Store implements storage.Storage interface.
// NB Store method works in async mode and only logs errors that have occured.
func (s *Storage) Store(orderUID, order string) error {
	s.storeCh <- storage.OrderDB{
		OrderUID:  orderUID,
		JSONOrder: order,
	}
	return nil
}

// storer is the worker function that listens to the storeCh channel and stores
// all incoming orders to the database.
func (s *Storage) storer() {
	log.Println("storage: postgres: storer started")
	for {
		select {
		case o := <-s.storeCh:
			if _, err := s.db.Exec(`INSERT INTO orders (uid, json_order) VALUES ($1, $2)`, o.OrderUID, o.JSONOrder); err != nil {
				log.Printf("storage: postgres: ERR: could not store the order %s: %s", o.OrderUID, err)
			} else {
				log.Printf("storage: postgres: order %s sucessfully stored", o.OrderUID)
			}
		case <-s.stopCh:
			log.Println("storage: postgres: storer stopped")
			s.wg.Done()
			return
		}
	}
}
