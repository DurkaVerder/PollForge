package service

import "question/models"

type Store interface {
	GetQuestions() ([]models.Question, error)
	UpdateCountAnswer(answer ...models.Answer) error
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{
		store: store,
	}
}
