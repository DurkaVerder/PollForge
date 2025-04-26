package store

import (
	"database/sql"
	"stream_line/internal/models"
)

const (
	GetOtherFormsQuery = `SELECT id, title, description, created_at, expires_at FROM forms WHERE creator_id != $1`
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) *Postgres {
	return &Postgres{
		db: db,
	}
}

func (p *Postgres) GetOtherForms(userID string) ([]models.FormFromDB, error) {
	rows, err := p.db.Query(GetOtherFormsQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var forms []models.FormFromDB
	for rows.Next() {
		var form models.FormFromDB
		if err := rows.Scan(&form.ID, &form.Title, &form.Description, &form.CreatedAt, &form.ExpiresAt); err != nil {
			return nil, err
		}
		forms = append(forms, form)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return forms, nil
}
