package emailnotifier

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
	
	return nil
}
