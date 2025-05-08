package storage

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func ConnectDB(driver, dns string) (*sql.DB, error) {
	db, err := sql.Open(driver, dns)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
