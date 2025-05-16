package service

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"profile/internal/models"
	"profile/internal/storage"
	"strings"
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
			&form.CreatorId,
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

func FormCheck(creatorId int, formId int) error {
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

func UpdateProfileName(userId int, profile models.UserProfile) error {
	err := storage.UpdateProfileNameRequest(userId, profile)
	if err != nil {
		log.Printf("Ошибка при обновлении профиля: %v", err)
	}
	return err
}

func DeleteProfile(userId int) error {
	err := storage.DeleteProfileRequest(userId)
	if err != nil {
		log.Printf("Ошибка при удалении профиля: %v", err)
	}
	return err
}

func UploadAvatar(userId int, avatarURL string) error {
	profile, err := storage.GetUserProfileRequest(userId)
	if err != nil {
		log.Printf("Ошибка при получении профиля: %v", err)
		return err
	}

	if profile.AvatarURL != "" {

		if !strings.HasPrefix(profile.AvatarURL, "/avatars/") {
			log.Printf("Некорректный путь аватара: %s", profile.AvatarURL)
			return fmt.Errorf("недопустимый путь аватара")
		}

		oldFilePath := strings.TrimPrefix(profile.AvatarURL, "/avatars/")
		fullPath := fmt.Sprintf("/uploads/avatars/%s", oldFilePath)
		log.Printf("Удаление старого аватара: %s", fullPath)
		if err := os.Remove(fullPath); err != nil {
			log.Printf("Ошибка при удалении старого аватара: %v", err)
		}
	}
	err = storage.UploadAvatarRequest(userId, avatarURL)
	if err != nil {
		log.Printf("Ошибка при загрузке аватара: %v", err)
		return fmt.Errorf("ошибка при загрузке аватара: %v", err)
	}
	return nil
}

func UpdateProfileBio(userId int, bio string) error {
	err := storage.UpdateProfileBioRequest(userId, bio)
	if err != nil {
		log.Printf("Ошибка при обновлении описании профиля: %v", err)
		return fmt.Errorf("Ошибка при обновлении описании профиля: %v", err)
	}
	return nil
}