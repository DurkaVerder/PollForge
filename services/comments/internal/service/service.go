package service

import (
	"comments/internal/models"
	"comments/internal/storage"
	"log"
)

func GetAllComments(formId int) ([]models.CommentResponse, error) {
	var comments []models.CommentResponse

	rows, err := storage.GetAllCommentsRequest(formId)
	if err != nil {
		log.Printf("Ошибка при получении данных комментариев: %v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var comment models.CommentResponse
		err := rows.Scan(&comment.Id, &comment.FormID, &comment.UserName, &comment.Description, &comment.CreatedAt, &comment.EditedAt)
		if err != nil {
			log.Printf("Ошибка при получении данных комментариев: %v", err)
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func CreateComment(comment models.CommentRequest, formId int, creatorId int) error {
	err := storage.CreateCommentRequest(comment, formId, creatorId)
	if err != nil {
		log.Printf("Ошибка при создании комментария: %v", err)
		return err
	}
	return nil
}

func UpdateUserComment(comment models.Comment, commentId int, formId int, creatorId int) error {
	err := storage.UpdateCommentRequest(comment, commentId, formId, creatorId)
	if err != nil {
		log.Printf("Ошибка при обновлении комментария: %v", err)
		return err
	}
	return nil
}

func DeleteComment(commentId int, formId int, creatorId int) error {
	err := storage.DeleteCommentRequest(commentId, formId, creatorId)
	if err != nil {
		log.Printf("Ошибка при удалении комментария: %v", err)
		return err
	}
	return nil
}
