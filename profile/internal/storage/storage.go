package storage

import (
	"database/sql"
	"log"
	"os"
	"profile/internal/models"

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
	err = Db.Ping()
	if err != nil {
		log.Printf("Ошибка доступа к базе данных: %v", err)
		return err
	}
	return nil
}

func GetUserProfileRequest(userId int) (*models.UserProfile, error) {
	row := Db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", userId)
	var profile models.UserProfile
	err := row.Scan(&profile.ID, &profile.Username, &profile.Email)
	if err != nil {
		log.Printf("Ошибка при получении профиля пользователя: %v", err)
		return nil, err
	}
	return &profile, nil
}

func GetUserFormsRequest(userId int) (*sql.Rows, error) {
	rows, err := Db.Query("SELECT id, title, description, link, private_key, expires_at FROM forms WHERE creator_id = $1", userId)
	if err != nil {
		log.Printf("Ошибка при получении форм пользователя: %v", err)
		return nil, err
	}
	return rows, nil
}

func FormCheckingRequest(existId int, creatorId int, formId int) error {

	queryCheck := "SELECT id FROM forms WHERE id  = $1 and creator_id = $2"
	err := Db.QueryRow(queryCheck, formId, creatorId).Scan(&existId)
	if err != nil {
		log.Printf("Ошибка при запросе на проверку наличия формы: %v", err)
		return err
	}
	return err
}

func FormDeleteRequest(formId int, creatorId int) error {

	query := "DELETE FROM forms WHERE id = $1 AND creator_id = $2"
	_, err := Db.Exec(query, formId, creatorId)
	if err != nil {
		log.Printf("Ошибка при запросе удаления формы: %v", err)
		return err
	}
	return err
}
