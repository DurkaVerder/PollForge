package service

import (
	"log"
	"math"
	"sort"
	"stream_line/internal/models"
)

type DB interface {
	GetFormsByOtherUserIDWithCountLikesAndComments(userID string) ([]models.FormFromDB, error)
	GetQuestionsByFormsID(formIDs []int) ([]models.QuestionFromDB, error)
	GetAnswersByQuestionsID(questionIDs []int, userID string) ([]models.AnswerFromDB, error)
}

type MergedData struct {
	Forms     []models.FormFromDB
	Questions []models.QuestionFromDB
	Answers   []models.AnswerFromDB
}

type FormWithQuestions struct {
	form      models.FormFromDB
	questions []QuestionWithAnswers
}

type QuestionWithAnswers struct {
	question models.QuestionFromDB
	answers  []models.AnswerFromDB
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

func (s *Service) GetStreamLines(userID string) (*models.StreamLineResponse, error) {
	forms, err := s.db.GetFormsByOtherUserIDWithCountLikesAndComments(userID)
	if err != nil {
		s.logger.Println("Error getting forms:", err)
		return nil, err
	}
	if len(forms) == 0 {
		s.logger.Println("No forms found")
		return nil, nil
	}
	formsIDs := s.getFormsIds(forms)
	questions, err := s.db.GetQuestionsByFormsID(formsIDs)
	if err != nil {
		s.logger.Println("Error getting questions:", err)
		return nil, err
	}

	questionIDs := s.getQuestionsIds(questions)
	answers, err := s.db.GetAnswersByQuestionsID(questionIDs, userID)
	if err != nil {
		s.logger.Println("Error getting answers:", err)
		return nil, err
	}
	s.logger.Printf("Forms: %v %v Questions: %v %v Answers: %v %v", len(forms), forms, len(questions), questions, len(answers), answers)
	
	mergedData := MergedData{
		Forms:     forms,
		Questions: questions,
		Answers:   answers,
	}

	polls := s.CreatePolls(mergedData)

	return &models.StreamLineResponse{Polls: polls}, nil
}

func (s *Service) CreatePolls(data MergedData) []models.Polls {

	if len(data.Forms) == 0 {
		return nil
	}

	questionWithAnswers := s.mergeQuestionWithAnswers(data.Questions, data.Answers)
	formsWithQuestions := s.mergeFormsWithQuestions(data.Forms, questionWithAnswers)

	polls := make([]models.Polls, 0, len(formsWithQuestions))
	for _, formWithQuestions := range formsWithQuestions {
		var poll models.Polls
		poll.ID = formWithQuestions.form.ID
		poll.Title = formWithQuestions.form.Title
		poll.Description = formWithQuestions.form.Description
		poll.Link = formWithQuestions.form.Link
		poll.Likes.Count = formWithQuestions.form.Like.Count
		poll.Likes.IsLiked = formWithQuestions.form.Like.IsLiked
		poll.CountVotes = formWithQuestions.form.CountVotes
		poll.CreatedAt = formWithQuestions.form.CreatedAt.Format("2006-01-02 15:04:05")
		poll.ExpiresAt = formWithQuestions.form.ExpiresAt.Format("2006-01-02 15:04:05")

		poll.Questions = make([]models.Question, 0, len(formWithQuestions.questions))
		for i, questionWithAnswers := range formWithQuestions.questions {
			poll.Questions[i].ID = questionWithAnswers.question.ID
			poll.Questions[i].Title = questionWithAnswers.question.Title
			poll.Questions[i].Answers = make([]models.Answer, len(questionWithAnswers.answers))

			for j, answer := range questionWithAnswers.answers {
				poll.Questions[i].Answers[j].ID = answer.ID
				poll.Questions[i].Answers[j].Title = answer.Title
				poll.Questions[i].Answers[j].Percent = answer.Percent
				poll.Questions[i].Answers[j].IsSelected = answer.IsSelected
			}
		}

		polls = append(polls, poll)
	}

	return polls
}

func (s *Service) mergeQuestionWithAnswers(questions []models.QuestionFromDB, answers []models.AnswerFromDB) []QuestionWithAnswers {
	questionMap := make(map[int][]models.AnswerFromDB)
	for _, answer := range answers {
		questionMap[answer.QuestionID] = append(questionMap[answer.QuestionID], answer)
	}

	result := make([]QuestionWithAnswers, 0, len(questions))
	for _, question := range questions {
		answers := questionMap[question.ID]
		sort.Slice(answers, func(i, j int) bool {
			return answers[i].NumberOrder < answers[j].NumberOrder
		})
		s.calculatePercent(answers)
		result = append(result, QuestionWithAnswers{
			question: question,
			answers:  answers,
		})
	}

	return result
}

func (s *Service) mergeFormsWithQuestions(forms []models.FormFromDB, questions []QuestionWithAnswers) []FormWithQuestions {
	formMap := make(map[int][]QuestionWithAnswers)
	for _, question := range questions {
		formMap[question.question.FormID] = append(formMap[question.question.FormID], question)
	}

	result := make([]FormWithQuestions, 0, len(forms))
	for _, form := range forms {
		questions := formMap[form.ID]
		sort.Slice(questions, func(i, j int) bool {
			return questions[i].question.NumberOrder < questions[j].question.NumberOrder
		})
		result = append(result, FormWithQuestions{
			form:      form,
			questions: questions,
		})
	}

	return result
}

func roundToTwoDecimalPlaces(val float64) float64 {
	return math.Round(val*100) / 100
}

func (s *Service) calculatePercent(answers []models.AnswerFromDB) {
	var totalVotes int
	for _, answer := range answers {
		totalVotes += answer.CountVotes
	}

	for i, answer := range answers {
		if totalVotes == 0 {
			answers[i].Percent = 0
		} else {
			percent := float64(answer.CountVotes*100) / float64(totalVotes)
			answers[i].Percent = roundToTwoDecimalPlaces(percent)
		}
	}
}

func (s *Service) getFormsIds(forms []models.FormFromDB) []int {
	ids := make([]int, len(forms))
	for i, form := range forms {
		ids[i] = form.ID
	}
	return ids
}

func (s *Service) getQuestionsIds(questions []models.QuestionFromDB) []int {
	ids := make([]int, len(questions))
	for i, question := range questions {
		ids[i] = question.ID
	}
	return ids
}
