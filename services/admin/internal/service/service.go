package service

import (
	"admin/internal/models"
	"admin/internal/storage"
	"fmt"
)

func GetAllUsers() ([]models.UserProfile, error) {
	rows, err := storage.GetAllUsersRequest()
	if err != nil {
		return nil, err
	}
	var users []models.UserProfile
	for rows.Next() {
		var user models.UserProfile
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Bio, &user.AvatarURL, &user.IsBanned); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("пользователи не найдены")
	}
	return users, nil
}

func ToggleBanUser(userId int, isBanned bool) error {
	err := storage.ToggleBanUserRequest(userId, isBanned)
	if err != nil {
		return fmt.Errorf("ошибка при изменении статуса пользователя: %v", err)
	}
	return nil
}

func DeleteUser(userId int) error {
	err := storage.DeleteUserRequest(userId)
	if err != nil {
		return fmt.Errorf("ошибка при удалении пользователя: %v", err)
	}
	return nil
}



func FormDelete(formId int) error {
	err := storage.FormDeleteRequest(formId)
	if err != nil {
		return fmt.Errorf("ошибка при удалении формы: %v", err)
	}
	return nil
}

func DeleteComment(commentId int) error {
	err := storage.DeleteCommentRequest(commentId)
	if err != nil {
		return fmt.Errorf("ошибка при удалении комментария: %v", err)
	}
	return nil
}