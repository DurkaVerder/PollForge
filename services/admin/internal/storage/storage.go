package storage

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var Db *sql.DB

func ConnectToDb() error {
	dsn := os.Getenv("DB_URL")
	var err error

	Db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Printf("Ошибка подключения к базе данных: %v", err)
		return err
	}
	Db.SetMaxOpenConns(50)
    Db.SetMaxIdleConns(25)
    Db.SetConnMaxIdleTime(5 * time.Minute)
	
	err = Db.Ping()
	if err != nil {
		log.Printf("Ошибка доступа к базе данных: %v", err)
		return err
	}
	return nil
}

func GetAllUsersRequest() (*sql.Rows, error) {
	rows, err := Db.Query("SELECT id, name, email, bio, avatar_url, is_banned FROM users")
	if err != nil {
		log.Printf("Ошибка при получении всех пользователей: %v", err)
		return nil, err
	}
	return rows, nil
}

func ToggleBanUserRequest(userId int, isBanned bool) error {
	query := "UPDATE users SET is_banned = $1 WHERE id = $2"
	_, err := Db.Exec(query, isBanned, userId)
	if err != nil {
		log.Printf("Ошибка при изменении статуса бана пользователя: %v", err)
		return err
	}
	return nil
}

func DeleteUserRequest(userId int) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := Db.Exec(query, userId)
	if err != nil {
		log.Printf("Ошибка при удалении пользователя: %v", err)
		return err
	}
	return nil
}

func FormDeleteRequest(formId int) error {
	query := "DELETE FROM forms WHERE id = $1"
	_, err := Db.Exec(query, formId)
	if err != nil {
		log.Printf("Ошибка при удалении формы: %v", err)
	}
	return err
}

func DeleteCommentRequest(commentId int) error {
	query := "DELETE FROM comments WHERE id = $1"
	_, err := Db.Exec(query, commentId)
	if err != nil {
		log.Printf("Ошибка при удалении комментария: %v", err)
		return err
	}
	return nil
}