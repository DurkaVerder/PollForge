package storage

import (
	"database/sql"
	"stream_line/internal/models"

	_ "github.com/lib/pq"
)

const (
	GetOtherFormsQuery = `SELECT f.id, f.title, f.description, l.count, EXISTS(SELECT * FROM likes_forms WHERE user_id = $1 AND form_id = f.id) AS is_liked, f.created_at, f.expires_at FROM forms f JOIN likes l ON l.form_id = f.id WHERE creator_id != $1`
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
		if err := rows.Scan(&form.ID, &form.Title, &form.Description, &form.Like.Count, &form.Like.IsLiked, &form.CreatedAt, &form.ExpiresAt); err != nil {
			return nil, err
		}
		forms = append(forms, form)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return forms, nil
}
