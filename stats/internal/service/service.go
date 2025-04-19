package service

import "stats/internal/models"

type Storage interface {
	GetQuestions(formID string) ([]models.QuestionFromDB, error)
	GetAnswers(formID string) ([]models.AnswerFromDB, error)
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) GetPollStats(userID, formID string) (models.PollStatsResponse, error) {

	questions, err := s.storage.GetQuestions(formID)
	if err != nil {
		return models.PollStatsResponse{}, err
	}

	answers, err := s.storage.GetAnswers(formID)
	if err != nil {
		return models.PollStatsResponse{}, err
	}

	pollStats := s.createQuestions(questions, answers)

	return models.PollStatsResponse{
		Stats: pollStats,
	}, nil
}

func (s *Service) createQuestions(questionsFromBD []models.QuestionFromDB, answersFromBD []models.AnswerFromDB) []models.PollStats {
	questions := make([]models.PollStats, len(questionsFromBD))

	for i, question := range questionsFromBD {
		questions[i].QuestionTitle = question.Title
		questions[i].Answers = make([]models.Answer, len(answersFromBD))

		for j, answer := range answersFromBD {
			questions[i].Answers[j] = models.Answer{
				Title: answer.Title,
				Count: answer.Count,
			}
		}
	}

	return questions
}
