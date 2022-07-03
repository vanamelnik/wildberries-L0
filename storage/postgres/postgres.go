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

type (
	Storage struct {
		db      *sql.DB
		storeCh chan orderDB
		stopCh  chan struct{}
		wg      *sync.WaitGroup
	}

	orderDB struct {
		orderUID  string
		jsonOrder string
	}
)

var _ storage.Storage = (*Storage)(nil)

//go:embed schema.sql
var queryCreate string

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
		storeCh: make(chan orderDB),
		stopCh:  make(chan struct{}),
		wg:      &sync.WaitGroup{},
	}
	s.wg.Add(1)
	go s.storer()
	return s, nil
}

func (s *Storage) Close() error {
	if s.stopCh != nil {
		close(s.stopCh)
		s.stopCh = nil
	}
	s.wg.Wait()
	return s.db.Close()
}

func (s *Storage) Get(orderUID string) (string, error) {
	var order string
	err := s.db.QueryRow(`SELECT order FROM orders WHERE uid = $1`, orderUID).Scan(&order)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrNotFound
		}
		return "", err
	}
	return order, nil
}

func (s *Storage) GetAll() ([]string, error) {
	orders := make([]string, 0)
	rows, err := s.db.Query("SELECT order FROM orders;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var order string
		if err := rows.Scan(&order); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (s *Storage) Store(orderUID, order string) error {
	s.storeCh <- orderDB{
		orderUID:  orderUID,
		jsonOrder: order,
	}
	return nil
}

func (s *Storage) storer() {
	log.Println("storage: postgres: storer started")
	for {
		select {
		case o := <-s.storeCh:
			if _, err := s.db.Exec(`INSERT INTO orders (uid, json_order) VALUES ($1, $2)`, o.orderUID, o.jsonOrder); err != nil {
				log.Printf("storage: postgres: ERR: could not store the order %s: %s", o.orderUID, err)
			} else {
				log.Printf("storage: postgres: order %s sucessfully stored", o.orderUID)
			}
		case <-s.stopCh:
			log.Println("storage: postgres: storer stopped")
			s.wg.Done()
			return
		}
	}
}
