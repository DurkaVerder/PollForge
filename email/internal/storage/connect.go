package storage

import "database/sql"

func ConnectDB(driver, dns string) *sql.DB {
	db, err := sql.Open(driver, dns)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}
	return db
}
