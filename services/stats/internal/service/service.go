package service

import (
	"stats/internal/models"

	"log"
)

type Storage interface {
	GetQuestions(formID string) ([]models.QuestionFromDB, error)
	GetAnswers(formID string) ([]models.AnswerFromDB, error)
	GetTimeChosen(answerIDs []int) ([]models.TimeChosenFromDB, error)

	GetProfileStats(userID string) (models.ProfileStatsFromDB, error)
	GetThemeStats(userID string) ([]models.Theme, error)
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

	answerIDs := s.GetAnswerIds(answers)
	timeChosen, err := s.storage.GetTimeChosen(answerIDs)
	if err != nil {
		return models.PollStatsResponse{}, err
	}
	s.mergeTimeChosenWithAnswers(answers, timeChosen)

	pollStats := s.createQuestions(questions, answers)

	return models.PollStatsResponse{
		Stats: pollStats,
	}, nil
}

func (s *Service) createQuestions(questionsFromBD []models.QuestionFromDB, answersFromBD []models.AnswerFromDB) []models.PollStats {
	questions := make([]models.PollStats, len(questionsFromBD))

	for i, question := range questionsFromBD {
		questions[i].QuestionTitle = question.Title
		questions[i].Answers = make([]models.Answer, 0, len(answersFromBD))

		for _, answer := range answersFromBD {
			if answer.QuestionID == question.ID {
				questions[i].Answers = append(questions[i].Answers, models.Answer{
					Title:        answer.Title,
					Count:        answer.Count,
					TimeSelected: answer.TimeSelected,
				})
			}
		}
	}

	return questions
}

func (s *Service) GetAnswerIds(answers []models.AnswerFromDB) []int {
	answerIds := make([]int, len(answers))
	for i, answer := range answers {
		answerIds[i] = answer.ID
	}
	return answerIds
}

func (s *Service) mergeTimeChosenWithAnswers(answers []models.AnswerFromDB, timeChosen []models.TimeChosenFromDB) {
	answerMap := make(map[int][]string)
	for _, tc := range timeChosen {
		answerMap[tc.IdAnswer] = append(answerMap[tc.IdAnswer], tc.Time.Format("2006-01-02 15:04:05"))
	}

	for i, answer := range answers {
		if times, exists := answerMap[answer.ID]; exists {
			answers[i].TimeSelected = times
		}
	}
}

func (s *Service) GetProfileStats(userID string) (models.ProfileStatsRequest, error) {
	baseInfo, err := s.storage.GetProfileStats(userID)
	if err != nil {
		log.Printf("error getting profile stats: %v", err)
		return models.ProfileStatsRequest{}, err
	}

	themes, err := s.storage.GetThemeStats(userID)
	if err != nil {
		log.Printf("error getting theme stats: %v", err)
		return models.ProfileStatsRequest{}, err
	}

	return models.ProfileStatsRequest{
		CountCreatedPolls:    baseInfo.CountCreated,
		CountAnsweredPolls:   baseInfo.CountAnswered,
		CountCommentsByPolls: baseInfo.CountComments,
		Themes:               themes,
	}, nil
}
