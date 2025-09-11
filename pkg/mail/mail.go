package mail

import (
	"api/internal/configuration"
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"

	"github.com/jordan-wright/email"
)

type MailService struct {
	host          string
	port          int
	fromName      string
	username      string
	password      string
	templatesPath string
}

func New(cfg *configuration.Mail) *MailService {
	return &MailService{
		host:          cfg.Host,
		port:          cfg.Port,
		fromName:      cfg.FromName,
		username:      cfg.Username,
		password:      cfg.Password,
		templatesPath: cfg.TemplatesPath,
	}
}

func (m *MailService) SendCode(to, purpose, code string) error {
	from := fmt.Sprintf("%s <%s>", m.fromName, m.username)
	data := map[string]any{
		"code": code,
	}

	templatePath := m.templatesPath + "/" + purpose + ".html"

	return m.sendWithTemplate(from, to, "This is Subjet", templatePath, data)
}

func (m *MailService) sendWithTemplate(from, to, subject, templatePath string, data interface{}) error {
	// load template
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("mail: failed to parse the template %s: %w", templatePath, err)
	}

	// set data to the template
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("mail: failed to execute template %s: %w", templatePath, err)
	}

	// setup mail body
	e := email.NewEmail()
	e.From = from
	e.To = []string{to}
	e.Subject = subject
	e.HTML = body.Bytes()

	addr := fmt.Sprintf("%s:%d", m.host, m.port)

	// send the message
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         m.host,
	}
	err = e.SendWithTLS(addr, LoginAuth(m.username, m.password), tlsConfig)
	if err != nil {
		return fmt.Errorf("mail: failed to send: %w", err)
	}

	return nil
}
