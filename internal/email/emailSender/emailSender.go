package emailSender

import (
	"bytes"
	"log"
	"os"

	"github.com/gophish/gomail"
)

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromAddress  string
}

type EmailSender struct {
	// Dialer handles the connection to SMTP server
	Config *EmailConfig
	Dialer *gomail.Dialer
}

func NewDefaultEmailSender() *EmailSender {
	emailConfig := EmailConfig{
		SMTPHost:     "smtp.gmail.com",
		SMTPPort:     587,
		SMTPUsername: os.Getenv("EMAIL_USERNAME"),
		SMTPPassword: os.Getenv("EMAIL_PASSWORD"),
		FromAddress:  os.Getenv("EMAIL_USERNAME"),
	}
	return NewEmailSender(&emailConfig)
}

func NewEmailSender(config *EmailConfig) *EmailSender {
	d := gomail.NewDialer(config.SMTPHost, config.SMTPPort, config.SMTPUsername, config.SMTPPassword)
	return &EmailSender{
		Config: config,
		Dialer: d,
	}
}

func (e *EmailSender) SendEmail(to string, subject string, body string, attachement *bytes.Buffer, attachmentName string, bodyType string) error {
	m := gomail.NewMessage()

	// Set the email headers
	m.SetHeader("From", e.Config.FromAddress)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody(bodyType, body)

	// Attach the required attachment
	if attachement != nil {
		m.AttachReader(attachmentName, attachement)
	}

	// Use the dialer to send the email
	if err := e.Dialer.DialAndSend(m); err != nil {
		log.Fatalf("Unable to send email. Error : %s", err)
		return err
	}

	return nil
}
