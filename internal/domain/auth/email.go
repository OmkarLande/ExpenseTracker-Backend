package domain_auth

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/smtp"
	"path/filepath"

	"ExpenseTracker-Backend/internal/config"
)

type EmailSender struct {
	Config *config.Config
}

func NewEmailSender(cfg *config.Config) *EmailSender {
	return &EmailSender{Config: cfg}
}

func (s *EmailSender) SendTemplateEmail(to string, subject string, templateName string, data interface{}) error {
	// Parse HTML template
	templatePath := filepath.Join("internal", "utils", "templates", templateName)
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("could not parse template: %w", err)
	}

	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s\n%s\n\n", subject, mimeHeaders)))

	err = tmpl.Execute(&body, data)
	if err != nil {
		return fmt.Errorf("could not execute template: %w", err)
	}

	// Setup Authentication
	auth := smtp.PlainAuth("", s.Config.SMTPUser, s.Config.SMTPPassword, s.Config.SMTPHost)

	// Send Email
	addr := fmt.Sprintf("%s:%s", s.Config.SMTPHost, s.Config.SMTPPort)

	if s.Config.SMTPPort == "465" {
		// Implicit TLS for port 465
		tlsConfig := &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         s.Config.SMTPHost,
		}

		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("could not dial TLS: %w", err)
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, s.Config.SMTPHost)
		if err != nil {
			return fmt.Errorf("could not create SMTP client: %w", err)
		}
		defer client.Close()

		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("could not authenticate SMTP client: %w", err)
		}

		if err = client.Mail(s.Config.SMTPFrom); err != nil {
			return fmt.Errorf("could not set SMTP sender: %w", err)
		}

		if err = client.Rcpt(to); err != nil {
			return fmt.Errorf("could not set SMTP recipient: %w", err)
		}

		w, err := client.Data()
		if err != nil {
			return fmt.Errorf("could not open SMTP data writer: %w", err)
		}
		_, err = w.Write(body.Bytes())
		if err != nil {
			return fmt.Errorf("could not write SMTP data: %w", err)
		}
		err = w.Close()
		if err != nil {
			return fmt.Errorf("could not close SMTP data writer: %w", err)
		}

		_ = client.Quit()
	} else {
		// STARTTLS or plaintext for other ports (e.g. 587)
		err = smtp.SendMail(addr, auth, s.Config.SMTPFrom, []string{to}, body.Bytes())
		if err != nil {
			return fmt.Errorf("could not send email: %w", err)
		}
	}

	return nil
}

func (s *EmailSender) SendVerificationEmail(to, token string) error {
    data := struct {
        Token string
    }{
        Token: token,
    }
    return s.SendTemplateEmail(to, "Verify your email", "verify_email.html", data)
}

func (s *EmailSender) SendPasswordResetEmail(to, otp string) error {
    data := struct {
        OTP string
    }{
        OTP: otp,
    }
    return s.SendTemplateEmail(to, "Password Reset OTP", "reset_password.html", data)
}