package storage

import (
	"auth/internal/models"
	"database/sql"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"

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

func Registration(hashedPassword []byte, request models.UserRequest) (string, error) {
	var userId string
	// Потому что мы пытаемся получить id после регистрации пользователя для создания jwt
	err := Db.QueryRow(`INSERT INTO users (name,email,password)
	                          VALUES($1,$2,$3) RETURNING id`,
		request.Username, request.Email, string(hashedPassword)).Scan(&userId)
	if err != nil {
		return "", err
	}
	return userId, err
}

func CheckingLoggingData(request models.UserRequest) (string, error) {
	var UserId string
	var hashedPassword []byte

	// Потому что если почта  совпадает, то мы получим id для генерации jwt и хешированный пароль для сравнения с обычным
	err := Db.QueryRow(`SELECT id, password FROM users WHERE email = $1`, request.Email).Scan(&UserId, &hashedPassword)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(request.Password))
	if err != nil {
		return "", err
	}
	return UserId, err

}
