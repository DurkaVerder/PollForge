package service

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"profile/internal/models"
	"profile/internal/storage"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func GetUserProfile(userId int) (*models.UserProfile, error) {
	profile, err := storage.GetUserProfileRequest(userId)
	if err != nil {
		log.Printf("Ошибка при получении профиля пользователя: %v", err)
		return nil, fmt.Errorf("пользователь не найден")
	}
	return profile, nil
}

func GetUserForms(userId int) ([]models.Form, error) {
	rows, err := storage.GetUserFormsRequest(userId)
	if err != nil {
		log.Printf("ошибка при получении форм в профиле: %v", err)
		return nil, err
	}
	defer rows.Close()

	var forms []models.Form
	for rows.Next() {
		var form models.Form
		err := rows.Scan(&form.Id,
			&form.Title,
			&form.Description,
			&form.Link,
			&form.PrivateKey,
			&form.ExpiresAt)
		if err != nil {
			log.Print("Данных нет")
			return nil, err
		}
		forms = append(forms, form)
	}

	return forms, nil
}

func GetToken(auth string) (*jwt.Token, error) {
	tokenStr := strings.TrimPrefix(auth, "Bearer")
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("неподходящий метод подписи")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	return token, err
}

func FormChek(creatorId int, formId int) error {
	var existId int
	err := storage.FormCheckingRequest(existId, creatorId, formId)
	if err != nil {
		log.Printf("Ошибка при проверке на наличие формы: %v", err)
		log.Printf("%s", err.Error())
		return err
	}
	return err
}

func FormDelete(formId int, creatorId int) (sql.Result, error) {
	err := storage.FormDeleteRequest(formId, creatorId)
	if err != nil {
		log.Printf("Ошибка при удалении данных: %v", err)
		return nil, err
	}
	return nil, err
}
