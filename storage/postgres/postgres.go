package postgres

import (
	"database/sql"
	_ "embed"
	"errors"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/vanamelnik/wildberries-L0/storage"
)

type (
	Storage struct {
		db      *sql.DB
		storeCh chan orderDB
		stopCh  chan struct{}
	}

	orderDB struct {
		orderUID  string
		jsonOrder string
	}
)

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
	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) Close() error {
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

}
