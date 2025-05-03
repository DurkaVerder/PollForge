package service

import (
	"database/sql"
	"fmt"
	"log"
	"profile/internal/models"
	"profile/internal/storage"

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



func FormChek(creatorId int, formId int) error {
	var existId int
	err := storage.FormCheckingRequest(existId, creatorId, formId)
	if err != nil {
		log.Printf("Ошибка при проверке на наличие формы: %v", err)
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
