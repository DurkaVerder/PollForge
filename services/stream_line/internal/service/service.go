package service

import (
	"log"
	"stream_line/internal/models"
)

type DB interface {
	GetOtherForms(userID string) ([]models.FormFromDB, error)
}

type Service struct {
	logger *log.Logger
	db     DB
}

func NewService(db DB, logger *log.Logger) *Service {
	return &Service{
		db:     db,
		logger: logger,
	}
}

func (s *Service) GetStreamLine(userID string) (models.FormResponse, error) {
	forms, err := s.db.GetOtherForms(userID)
	if err != nil {
		s.logger.Printf("Error getting forms from DB: %v", err)
		return models.FormResponse{}, err
	}

	response := s.createFormResponse(forms)

	return response, nil
}

func (s *Service) createFormResponse(forms []models.FormFromDB) models.FormResponse {
	response := models.FormResponse{
		Forms: make([]models.Form, len(forms)),
	}

	for i, form := range forms {
		response.Forms[i] = models.Form{
			ID:          form.ID,
			Title:       form.Title,
			Description: form.Description,
			Likes: models.Like{
				Count:   form.Like.Count,
				IsLiked: form.Like.IsLiked,
			},
			CreatedAt: form.CreatedAt.Format("2006-01-02 15:04:05"),
			ExpiresAt: form.ExpiresAt.Format("2006-01-02 15:04:05"),
		}
	}

	return response
}
