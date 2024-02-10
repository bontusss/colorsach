package utils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/bontusss/colosach/config"
	"github.com/bontusss/colosach/models"
	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

type EmailData struct {
	URL       string
	FirstName string
	Subject   string
}

// ParseTemplateDir 👇 Email template parser
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

// SendEmail ? Email template parser
func SendEmail(user *models.DBResponse, data *EmailData, templateName string) error {
	loadConfig, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("could not load loadConfig", err)
	}

	// Sender data.
	from := loadConfig.EmailFrom
	smtpPass := loadConfig.SMTPPass
	smtpUser := loadConfig.SMTPUser
	to := user.Email
	smtpHost := loadConfig.SMTPHost
	smtpPort := loadConfig.SMTPPort

	var body bytes.Buffer

	templateDir, err := ParseTemplateDir("templates")
	if err != nil {
		log.Fatal("Could not parse templateDir", err)
	}

	templateDir = templateDir.Lookup(templateName)
	err = templateDir.Execute(&body, &data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(templateDir.Name())

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil

}