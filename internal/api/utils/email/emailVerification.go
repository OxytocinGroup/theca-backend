package utils

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/OxytocinGroup/theca-backend/internal/config"
	"github.com/resend/resend-go/v2"
)

type Mail struct {
	Email    string
	Username string
	Code     string
}

func (m *Mail) SendVerificationEmail(cfg *config.Config, email, code, username string) error {
	template := template.New("mail.html")

	template, err := template.ParseFiles("internal/api/utils/email/mail.html")
	if err != nil {
		return err
	}

	var tpl bytes.Buffer
	if err := template.Execute(&tpl, Mail{Username: username, Code: code}); err != nil {
		return err
	}

	apiKey := cfg.SMTPAPI

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "Theca <no-reply@theca.oxytocingroup.com>",
		To:      []string{email},
		Html:    tpl.String(),
		Subject: fmt.Sprintf("%s | Verification Code", code),
	}

	_, pohuy1 := client.Emails.Send(params)
	return pohuy1
}
