package service

import (
	"log"
	"question/models"
	"strconv"
	"sync"
)

type Store interface {
	GetQuestions(formID int) ([]models.QuestionFromDB, error)
	GetAnswers(formID int) ([]models.AnswerFromDB, error)
	UpdateCountAnswer(ids []int) error
}

type Service struct {
	store          Store
	answersChannel chan []models.SubmitAnswer
}

func NewService(store Store, answerChannel chan []models.SubmitAnswer) *Service {
	return &Service{
		store:          store,
		answersChannel: answerChannel,
	}
}

func (s *Service) GetQuestions(formID string) (models.QuestionResponse, error) {
	id, err := strconv.Atoi(formID)
	if err != nil {
		log.Printf("GetQuestions: Ошибка при преобразовании formID: %v\n", err)
		return models.QuestionResponse{}, err
	}

	questionsFromBD, err := s.store.GetQuestions(id)
	if err != nil {
		log.Printf("GetQuestions: Ошибка при получении вопросов: %v\n", err)
		return models.QuestionResponse{}, err
	}

	answersFromBD, err := s.store.GetAnswers(id)
	if err != nil {
		log.Printf("GetQuestions: Ошибка при получении ответов: %v\n", err)
		return models.QuestionResponse{}, err
	}

	questions := s.createQuestions(questionsFromBD, answersFromBD)
	questionsResponse := models.QuestionResponse{
		Question: questions,
	}

	return questionsResponse, nil
}

func (s *Service) AddAnswerRequestToChannel(answer models.SubmitAnswerRequest) {
	s.answersChannel <- answer.Answers
}

func (s *Service) createQuestions(question []models.QuestionFromDB, answers []models.AnswerFromDB) []models.Question {
	questions := make([]models.Question, 0, len(question))

	for _, q := range question {
		question := models.Question{
			ID:    q.ID,
			Title: q.Title,
		}

		for _, a := range answers {
			if a.QuestionID == q.ID {
				question.Answers = append(question.Answers, models.Answer{
					ID:    a.ID,
					Title: a.Title,
				})
			}
		}

		questions = append(questions, question)
	}

	return questions
}

func (s *Service) writeAnswer(answers []models.SubmitAnswer) error {
	ids := s.getSelectedIds(answers)

	if len(ids) == 0 {
		log.Println("WriteAnswer: Нет выбранных ответов")
		return nil
	}

	err := s.store.UpdateCountAnswer(ids)
	if err != nil {
		log.Printf("WriteAnswer: Ошибка при обновлении счетчика ответов: %v\n", err)
		return err
	}

	return nil
}

func (s *Service) getSelectedIds(answers []models.SubmitAnswer) []int {
	idsSelectedAnswer := make([]int, 0, len(answers)/2)

	for _, answer := range answers {
		if answer.Selected {
			idsSelectedAnswer = append(idsSelectedAnswer, answer.ID)
		}
	}

	return idsSelectedAnswer
}

func (s *Service) writeAnswerWorker(wg *sync.WaitGroup) {
	defer wg.Done()
	for val := range s.answersChannel {
		err := s.writeAnswer(val)
		if err != nil {
			log.Printf("WriteAnswerWorker: Ошибка при записи ответа: %v\n", err)
			continue
		}
		log.Printf("WriteAnswerWorker: Ответ записан: %+v\n", val)
	}
}

func (s *Service) Close() {
	close(s.answersChannel)
}

func (s *Service) StartWorker(countWorkers int, wg *sync.WaitGroup) {
	if countWorkers <= 0 {
		panic("Количество воркеров должно быть больше 0")
	}

	for i := 0; i < countWorkers; i++ {
		wg.Add(1)
		go s.writeAnswerWorker(wg)
	}
}
