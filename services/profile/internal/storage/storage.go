package storage

import (
	"database/sql"
	"log"
	"os"
	"profile/internal/models"

	_ "github.com/lib/pq" // PostgreSQL driver
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
	row := Db.QueryRow("SELECT id, name, email, COALESCE(bio, ''), COALESCE(avatar_url, '') FROM users WHERE id = $1", userId)
	var profile models.UserProfile
	err := row.Scan(&profile.ID, &profile.Username, &profile.Email, &profile.Bio, &profile.AvatarURL)
	if err != nil {
		log.Printf("Ошибка при получении профиля пользователя: %v", err)
		return nil, err
	}
	return &profile, nil
}

func GetUserFormsRequest(userId int) (*sql.Rows, error) {
	rows, err := Db.Query("SELECT id, creator_id, title, description, link, private_key, expires_at FROM forms WHERE creator_id = $1", userId)
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

func UpdateProfileNameRequest(userId int, profile models.UserProfile) error {
	query := "UPDATE users SET name = $1 WHERE id = $2"
	_, err := Db.Exec(query, profile.Username, userId)
	if err != nil {
		log.Printf("Ошибка при обновлении профиля пользователя: %v", err)
		return err
	}
	return nil
}

func DeleteProfileRequest(userId int) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := Db.Exec(query, userId)
	if err != nil {
		log.Printf("Ошибка при удалении профиля пользователя: %v", err)
		return err
	}
	return nil
}

func UploadAvatarRequest(userId int, avatarURL string) error {
	query := "UPDATE users SET avatar_url = $1 WHERE id = $2"
	_, err := Db.Exec(query, avatarURL, userId)
	if err != nil {
		log.Printf("Ошибка при загрузке аватара: %v", err)
		return err
	}
	return nil
}

func UpdateProfileBioRequest(userId int, bio string) error {
	query := "UPDATE users SET bio = $1 WHERE id = $2"
	_, err := Db.Exec(query, bio, userId)
	if err != nil {
		log.Printf("Ошибка при обновлении профиля пользователя: %v", err)
		return err
	}
	return nil
}
