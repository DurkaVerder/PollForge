package service

import (
	en "email/internal/email_notifier"
	"email/internal/models"
	"fmt"
	"log"
	"os"
	"sync"
)

const (
	countWorker = 5

	userRegisteredEvent = "user_registered"

	userLoginEvent      = "user_login"

	userPasswordEvent 	= "user_reset"

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

		emailMsg := s.createEmail(email, m.EventType, m.Token)

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

func (s *Service) selectEmailTemplate(eventType, token string) (subject, body string) {
    switch eventType {
    case userRegisteredEvent:
        subject = "Добро пожаловать в наш сервис!"
        body = "Спасибо за регистрацию! Мы рады видеть вас."

    case userLoginEvent:
        subject = "Уведомление о входе в систему"
        body = "Вы успешно вошли. Если это были не вы — срочно смените пароль!"

    case userPasswordEvent:
        subject = "Сброс пароля"
        base := os.Getenv("FRONTEND_URL") 
        link := fmt.Sprintf("%s/reset_password?token=%s", base, token)
        body = fmt.Sprintf("Чтобы сбросить пароль, перейдите по ссылке:\n\n%s\n\nЕсли вы не запрашивали сброс — проигнорируйте это письмо.", link)

    default:
        subject = "У вас новое уведомление"
        body    = "Пожалуйста, проверьте ваш аккаунт."
    }

    return subject, body
}

func (s *Service) createEmail(email, eventType string, token string) Email {

	subject, body := s.selectEmailTemplate(eventType, token)

	return Email{
		To:      email,
		Subject: subject,
		Body:    body,
	}
}

