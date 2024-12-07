package utils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"text/template"

	"github.com/OxytocinGroup/theca-backend/internal/config"
	"gopkg.in/gomail.v2"
)

type Mail struct {
	Email    string
	Username string
	Code     string
}

func (m *Mail) SendVerificationEmail(cfg config.Config, email, code, username string) error {
	t := template.New("mail.html")

	t, err := t.ParseFiles("internal/api/utils/email/mail.html")
	if err != nil {
		return err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, Mail{Username: username, Code: code}); err != nil {
		return err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", cfg.SMTPFrom)
	message.SetHeader("To", email)
	message.SetHeader("Subject", "Email verification")
	message.SetBody("text/html", tpl.String())

	fmt.Println(
		"DEBUUUUG:",
		cfg.SMTPServer,
		cfg.SMTPPort,
		cfg.SMTPUsername,
		cfg.SMTPPassword,
		cfg.SMTPFrom,
	)
	d := gomail.NewDialer(
		cfg.SMTPServer,
		cfg.SMTPPort,
		cfg.SMTPUsername,
		cfg.SMTPPassword,
	)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(message); err != nil {
		return err
	}
	return nil
}
