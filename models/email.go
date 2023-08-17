package models

import "github.com/go-mail/mail/v2"

const (
	DefaultSender = "support@webapp.com"
)

type Email struct {
	//
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type EmailService struct {
	// DefaultSender is used when one is not provided,
	// this is also used in functions where the email
	// is predetermined, i.e. the forgot password email.
	DefaultSender string

	// unexported
	dialer *mail.Dialer
}

func NewEmailService(config SMTPConfig) (*EmailService, error) {
	es := EmailService{
		dialer: mail.NewDialer(config.Host, config.Port, config.Username, config.Password),
	}
	return &es, nil
}
