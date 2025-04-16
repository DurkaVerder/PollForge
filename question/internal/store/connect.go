package store

import (
	"database/sql"
	"os"
	"time"
)

const (
	maxOpenConns    = 25
	maxIdleConns    = 25
	connMaxLifetime = 1 * time.Minute
)

func InitConnection() (*sql.DB, error) {
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)

	return db, nil
}
