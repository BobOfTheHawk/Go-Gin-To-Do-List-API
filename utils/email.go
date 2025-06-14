package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

// SendEmail sends an email using the SMTP settings from environment variables.
func SendEmail(to, subject, body string) error {
	// Get SMTP configuration from environment variables
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASSWORD")

	// If any of the required SMTP variables are not set, log an error and return.
	// This prevents the app from crashing and allows it to run without email functionality.
	if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPass == "" {
		return fmt.Errorf("SMTP configuration is incomplete. Email not sent")
	}

	// The `auth` variable is used to authenticate with the SMTP server.
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	// The message is formatted according to RFC 822.
	msg := "From: " + smtpUser + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	// `SendMail` connects to the server, authenticates, and sends the email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{to}, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Successfully sent email to %s", to)
	return nil
}
