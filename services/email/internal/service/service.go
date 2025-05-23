package service

import (
	en "email/internal/email_notifier"
	"email/internal/models"
	"log"
	"sync"
)

const (
	countWorker = 5

	userRegisteredEvent = "user_registered"

	userLoginEvent      = "user_login"

)

type DB interface {
	GetEmailByUserID(userID string) (string, error)
}

type Email struct {
	To      string
	Subject string
	Body    string
}

type Service struct {
	wg     *sync.WaitGroup
	logger *log.Logger
	db     DB

	emailNotifier *en.EmailNotifier
}

func NewService(db DB, emailNotifier *en.EmailNotifier, logger *log.Logger) *Service {
	return &Service{
		wg:            &sync.WaitGroup{},
		logger:        logger,
		db:            db,
		emailNotifier: emailNotifier,
	}
}

func (s *Service) StartWorker(msg <-chan models.MessageKafka) {
	for i := 0; i < countWorker; i++ {
		s.wg.Add(1)
		go s.worker(msg)
	}
}

func (s *Service) StopWorker() {
	s.wg.Wait()
}

func (s *Service) worker(msg <-chan models.MessageKafka) {
	defer s.wg.Done()
	for m := range msg {
		email, err := s.getEmailByUserID(m.UserID)
		if err != nil {
			s.logger.Println("Ошибка получения почты по id пользователя:", err)
			continue
		}

		emailMsg := s.createEmail(email, m.EventType)

		err = s.emailNotifier.SendEmail(emailMsg.To, emailMsg.Subject, emailMsg.Body)
		if err != nil {
			s.logger.Println("Ошибка отправки письма:", err)
			continue
		}
	}
}

func (s *Service) getEmailByUserID(userID string) (string, error) {
	email, err := s.db.GetEmailByUserID(userID)
	if err != nil {
		return "", err
	}
	return email, nil
}

func (s *Service) selectEmailTemplate(eventType string) (string, string) {
	var subject, body string
	switch eventType {
	case userRegisteredEvent:

		subject = "Добро пожаловать в наш сервис!"
		body = "Спасибо за регистрацию! Мы рады видеть вас в нашем сервисе."
	case userLoginEvent:
		subject = "Уведомление о входе в систему!"
		body = "Вы успешно вошли в систему. Если это были не вы, пожалуйста, измените пароль."

	default:
		subject = "Уведомление"
		body = "У вас новое уведомление."
	}

	return subject, body
}

func (s *Service) createEmail(email, eventType string) Email {

	subject, body := s.selectEmailTemplate(eventType)

	return Email{
		To:      email,
		Subject: subject,
		Body:    body,
	}
}
