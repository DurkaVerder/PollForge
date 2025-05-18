package storage

import (
	"auth/internal/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

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
func GetUserIDByEmail(email string) (string, error) {
    var id string
    err := Db.QueryRow(`SELECT id FROM users WHERE email = $1`, email).Scan(&id)
    if err == sql.ErrNoRows {
        return "", errors.New("пользователь не найден")
    }
    if err != nil {
        return "", fmt.Errorf("GetUserIDByEmail: %w", err)
    }
    return id, nil
}

func CreatePasswordReset(userID string, token string, expiresAt time.Time) error {
    query := `
      INSERT INTO password_resets (user_id, token, expires_at)
      VALUES ($1, $2, $3)
    `
    if _, err := Db.Exec(query, userID, token, expiresAt); err != nil {
        return fmt.Errorf("CreatePasswordReset: %w", err)
    }
    return nil
}

func GetPasswordResetByToken(token string) (*models.PasswordReset, error) {
    const query = `
        SELECT id, user_id, token, expires_at
        FROM password_resets
        WHERE token = $1
    `
    pr := &models.PasswordReset{}
    err := Db.QueryRow(query, token).
        Scan(&pr.ID, &pr.UserID, &pr.Token, &pr.ExpiresAt)
    if err == sql.ErrNoRows {
        return nil, errors.New("токен сброса не найден")
    }
    if err != nil {
        return nil, fmt.Errorf("GetPasswordResetByToken: %w", err)
    }
    return pr, nil
}

func DeletePasswordReset(id int) error {
    const query = `
        DELETE FROM password_resets
        WHERE id = $1
    `
    res, err := Db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("DeletePasswordReset: %w", err)
    }
    rows, err := res.RowsAffected()
    if err != nil {
        return fmt.Errorf("DeletePasswordReset.RowsAffected: %w", err)
    }
    if rows == 0 {
        return errors.New("запись сброса не найдена для удаления")
    }
    return nil
}

func UpdateUserPassword(userID int, hashedPassword string) error {
    const query = `
        UPDATE users
        SET password = $1
        WHERE id = $2
    `
    res, err := Db.Exec(query, hashedPassword, userID)
    if err != nil {
        return fmt.Errorf("UpdateUserPassword: %w", err)
    }
    rows, err := res.RowsAffected()
    if err != nil {
        return fmt.Errorf("UpdateUserPassword.RowsAffected: %w", err)
    }
    if rows == 0 {
        return errors.New("пользователь не найден при обновлении пароля")
    }
    return nil
}