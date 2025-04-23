package storage

import "database/sql"

const (
	GetEmailByUserIDQuery = "SELECT email FROM users WHERE id = $1"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) *Postgres {
	return &Postgres{db: db}
}

func (p *Postgres) GetEmailByUserID(userID string) (string, error) {
	var email string
	err := p.db.QueryRow(GetEmailByUserIDQuery, userID).Scan(&email)
	if err != nil {
		return "", err
	}
	return email, nil
}
