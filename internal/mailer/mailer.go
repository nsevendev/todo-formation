package mailer

import (
	"fmt"
	"net/smtp"
	"os"
)

type MailData struct {
	To      string
	Subject string
	Body    string
}

func SendMail(m MailData) error {
	// TODO : changer les variables d'environnement avec le pkg config
	auth := smtp.PlainAuth(
		"",
		os.Getenv("MAILTRAP_USER"),
		os.Getenv("MAILTRAP_PASS"),
		os.Getenv("MAILTRAP_HOST"),
	)

	from := os.Getenv("MAIL_FROM")
	to := []string{m.To}

	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		from, m.To, m.Subject, m.Body,
	))

	addr := fmt.Sprintf("%s:%s", os.Getenv("MAILTRAP_HOST"), os.Getenv("MAILTRAP_PORT"))
	return smtp.SendMail(addr, auth, from, to, msg)
}
