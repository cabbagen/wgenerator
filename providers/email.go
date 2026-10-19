package providers

import (
	"gopkg.in/gomail.v2"
)

type Email struct {
	From    string
	To      []string
	Cc      []string
	Subject string
	Body    map[string]string
}

type EmailSender struct {
	Host     string
	Port     int
	Username string
	Password string
}

func (sender EmailSender) HandleSendEmail(email Email) error {
	m, headersMap := gomail.NewMessage(), map[string][]string{
		"From":    {email.From},
		"To":      email.To,
		"Cc":      email.Cc,
		"Subject": {email.Subject},
	}

	m.SetHeaders(headersMap)

	for key, value := range email.Body {
		m.SetBody(key, value)
	}

	d := gomail.NewDialer(sender.Host, sender.Port, sender.Username, sender.Password)

	return d.DialAndSend(m)
}
