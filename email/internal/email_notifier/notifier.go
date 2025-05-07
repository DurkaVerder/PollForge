package emailnotifier

import (
	"log"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type EmailNotifier struct {
	from     string
	password string
	smtpHost string
	smtpPort string
}

func NewEmailNotifier(from, password, smtpHost, smtpPort string) *EmailNotifier {
	return &EmailNotifier{
		from:     from,
		password: password,
		smtpHost: smtpHost,
		smtpPort: smtpPort,
	}
}

func (emailNotifier *EmailNotifier) SendEmail(to, subject, body string) error {
	log.Println("Отправка на почту:", to)
	e := email.NewEmail()
	e.From = emailNotifier.from
	e.To = []string{to}
	e.Subject = subject
	e.Text = []byte(body)
	err := e.Send(emailNotifier.smtpHost+":"+emailNotifier.smtpPort, smtp.PlainAuth("", emailNotifier.from, emailNotifier.password, emailNotifier.smtpHost))
	return err
}
