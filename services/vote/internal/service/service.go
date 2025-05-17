package service

import (
	"log"
	"question/models"
	"sync"
)

type Storage interface {
	UpdateCountAnswerUp(answerId, userID int) error
	UpdateCountAnswerDown(answerId, userID int) error
}

type Service struct {
	wg             *sync.WaitGroup
	storage        Storage
	answersChannel chan models.Vote
}

func NewService(storage Storage, answerChannel chan models.Vote) *Service {
	return &Service{
		wg:             &sync.WaitGroup{},
		storage:        storage,
		answersChannel: answerChannel,
	}
}

func (s *Service) AddVoteRequestToChannel(vote models.Vote) {
	s.answersChannel <- vote
}

func (s *Service) UpdateCountAnswer(vote models.Vote) error {
	if vote.IsUpVote {
		return s.storage.UpdateCountAnswerUp(vote.ID, vote.UserID)
	} else {
		return s.storage.UpdateCountAnswerDown(vote.ID, vote.UserID)
	}
}

func (s *Service) ProcessVote() {
	defer s.wg.Done()
	for vote := range s.answersChannel {
		if err := s.UpdateCountAnswer(vote); err != nil {
			log.Printf("ProcessVote: Ошибка при обновлении количества голосов: %v\n", err)
			continue
		}
		log.Printf("ProcessVote: Успешно обновлено количество голосов для ответа ID %d, пользователь ID %d\n", vote.ID, vote.UserID)
	}
}

func (s *Service) Start(countWorker int) {
	for i := 0; i < countWorker; i++ {
		s.wg.Add(1)
		go s.ProcessVote()
	}
}

func (s *Service) Close() {
	close(s.answersChannel)
	s.wg.Wait()
}
