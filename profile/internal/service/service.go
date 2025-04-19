package service

import (
	"fmt"
	"log"
	"profile/internal/models"
	"profile/internal/storage"
)

func GetUserProfile(userId int) (*models.UserProfile, error) {
	row := storage.Db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", userId)
	var profile models.UserProfile
	err := row.Scan(&profile.ID, profile.Username, profile.Email)
	if err != nil {
		return nil, fmt.Errorf("Пользователь не найден")
	}
	return &profile, nil
}

func GetUserForms(userId int) ([]models.Form, error) {
	rows, err := storage.Db.Query("SELECT id, title, description, link FROM forms WHERE creator_id = $1", userId)
	if err != nil {
		log.Print("Данных нет")
		return nil, err
	}
	var forms []models.Form

	defer rows.Close()

	for rows.Next() {
		var f models.Form
		err := rows.Scan(&f.ID, &f.Title, &f.Description, &f.Link)
		if err != nil {
			log.Print("Данных нет")
			return nil, err
		}
		forms = append(forms, f)
	}

	return forms, nil
}
