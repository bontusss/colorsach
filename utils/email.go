package utils

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/bontusss/colosach/models"
	"github.com/wneessen/go-mail"

	"github.com/k3a/html2text"
)

type EmailData struct {
	URL       string
	Username string
	Subject   string
}

// 👇 Email template parser

func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	fmt.Println("Am parsing templates...")

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}

func SendEmail(user *models.DBResponse, data *EmailData, templateName string) error {
	// Sender data.
	from := os.Getenv("EMAIL_FROM")
	// from := "colosach.app@gmail.com"
	smtpPass := os.Getenv("SMTP_PASS")
	// smtpPass := "kfbk fiop hbry zmom"
	smtpUser := os.Getenv("SMTP_USER")
	to := user.Email
	smtpHost := os.Getenv("SMTP_HOST")
	// smtpHost := "smtp.gmail.com"
	smtpPort := os.Getenv("SMTP_PORT")
	port, err := strconv.Atoi(smtpPort)
	if err != nil {
		return err
	}

	var body bytes.Buffer

	template, err := ParseTemplateDir("templates")
	if err != nil {
		log.Fatal("Could not parse template", err)
	}

	template = template.Lookup(templateName)
	template.Execute(&body, &data)
	fmt.Println(template.Name())

	m := mail.NewMsg()
	if err := m.From(from); err != nil {
		log.Fatalf("failed to set From address: %s", err)
	}
	if err := m.To(to); err != nil {
		log.Fatalf("failed to set To address: %s", err)
	}
	m.Subject(data.Subject)
	m.SetBodyString(mail.TypeTextPlain, html2text.HTML2Text(body.String()))
	c, err := mail.NewClient(smtpHost, mail.WithPort(port), mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithUsername(smtpUser), mail.WithPassword(smtpPass))
	if err != nil {
		log.Fatalf("failed to create mail client: %s", err)
	}
	if err := c.DialAndSend(m); err != nil {
		log.Fatalf("failed to send mail: %s", err)
	}
	return nil
}
