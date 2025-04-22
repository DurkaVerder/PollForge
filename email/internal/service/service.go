package service

import (
	en "email/internal/email_notifier"
	"email/internal/models"
	"sync"
)

const (
	countWorker = 5

	userRegisteredTemplate = "user_registered"
	userLoginTemplate      = "user_login"
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
	wg *sync.WaitGroup
	db DB

	emailNotifier en.EmailNotifier
}

func NewService(db DB, emailNotifier en.EmailNotifier) *Service {
	return &Service{
		wg:            &sync.WaitGroup{},
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
		email, err := s.createEmail(m)
		if err != nil {
			continue
		}

		err = s.emailNotifier.SendEmail(email.To, email.Subject, email.Body)
		if err != nil {
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
	case userRegisteredTemplate:
		subject = "Welcome to our service!"
		body = "Thank you for registering. We are glad to have you."
	case userLoginTemplate:
		subject = "Login Notification"
		body = "You have successfully logged in to your account."
	}
	return subject, body
}

func (s *Service) createEmail(msg models.MessageKafka) (Email, error) {
	email, err := s.getEmailByUserID(msg.UserID)
	if err != nil {
		return Email{}, err
	}

	subject, body := s.selectEmailTemplate(msg.EventType)

	return Email{
		To:      email,
		Subject: subject,
		Body:    body,
	}, nil
}
