package service

import (
	"log"
	"question/internal/models"
	"sync"
)

type Storage interface {
	UpdateCountAnswerUp(answerId, userID int) error
	UpdateCountAnswerDown(answerId, userID int) error

	UpdateCountLikeUp(formID, userID int) error
	UpdateCountLikeDown(formID, userID int) error
}

type Service struct {
	wg             *sync.WaitGroup
	storage        Storage
	answersChannel chan models.Vote
	likesChannel   chan models.Like
}

func NewService(storage Storage, answerChannel chan models.Vote, likesChannel chan models.Like) *Service {
	return &Service{
		wg:             &sync.WaitGroup{},
		storage:        storage,
		answersChannel: answerChannel,
		likesChannel:   likesChannel,
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
		s.wg.Add(1)
		go s.ProcessLike()
	}
}

func (s *Service) Close() {
	close(s.answersChannel)
	close(s.likesChannel)
	s.wg.Wait()
}

func (s *Service) AddLikeRequestToChannel(like models.Like) {
	s.likesChannel <- like
}

func (s *Service) UpdateCountLike(like models.Like) error {
	if like.IsUpLike {
		return s.storage.UpdateCountLikeUp(like.ID, like.UserID)
	} else {
		return s.storage.UpdateCountLikeDown(like.ID, like.UserID)
	}
}

func (s *Service) ProcessLike() {
	defer s.wg.Done()
	for like := range s.likesChannel {
		if err := s.UpdateCountLike(like); err != nil {
			log.Printf("ProcessLike: Ошибка при обновлении количества Like: %v\n", err)
			continue
		}
		log.Printf("ProcessLike: Успешно обновлено количество Like для формы ID %d, пользователь ID %d\n", like.ID, like.UserID)
	}
}
