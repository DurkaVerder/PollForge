package storage

import (
	"database/sql"
	"time"
)

const (
	maxOpenConns    = 25
	maxIdleConns    = 25
	connMaxLifetime = 1 * time.Minute
)

func InitConnection(driverName, connURL string) (*sql.DB, error) {
	db, err := sql.Open(driverName, connURL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)

	return db, nil
}
