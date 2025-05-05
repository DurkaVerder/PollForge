package storage

import (
	"comments/internal/models"
	"database/sql"
	"log"
	"os"
)

var Db *sql.DB

func ConnectToDb() error {
	dsn := os.Getenv("DB_URL")
	var err error

	Db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Printf("Ошибка подключения к базе данных: %v", err)
	}

	err = Db.Ping()

	if err != nil {
		log.Printf("Ошибка доступа к базе данных: %v", err)
	}
	return nil
}

func GetAllCommentsRequest(formId int)(*sql.Rows, error) {
	query := `SELECT description, created_at, updated_at FROM comments WHERE form_id = $1 ORDER BY created_at DESC`

	rows, err := Db.Query(query, formId)

	if err != nil {
		log.Printf("Ошибка при запросе получения всех комментариев: %v", err)
		return nil, err
	}
	return rows, nil
}

func CreateCommentRequest(comment models.Comment, formId int, creatorId int) error {
	query := `INSERT INTO comments (form_id, user_id, description, created_at) VALUES ($1, $2, $3, $4)`
	_, err := Db.Exec(query, formId, creatorId, comment.Description, comment.CreatedAt)
	if err != nil {
		log.Printf("Ошибка при запросе создания комментария: %v", err)
		return err
	}
	return nil
}

func UpdateCommentRequest(comment models.Comment, commentId int, formId int, creatorId int) error {
	query := `UPDATE comments SET description = $1, WHERE form_id = $2 AND id = $3 AND user_id = $4`
	_, err := Db.Exec(query, comment.Description, formId, commentId, creatorId)
	if err != nil {
		log.Printf("Ошибка при запросе обновления комментария: %v", err)
		return err
	}
	return nil
}

func DeleteCommentRequest(commentId int, formId int, creatorId int) error {
	query := `DELETE FROM comments WHERE form_id = $1 AND id = $2 AND user_id = $3`
	_, err := Db.Exec(query, formId, commentId, creatorId)
	if err != nil {
		log.Printf("Ошибка при запросе удаления комментария: %v", err)
		return err
	}
	return nil
}