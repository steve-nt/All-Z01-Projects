package utils

import (
	"fmt"
	"net"
	"net/smtp"
	"time"
)

func SendVerificationEmail(to, link string) error {
	from := "jay23061984@gmail.com"
	password := "huvx apns smrh smrq" // App Password από Gmail

	subject := "Subject: Confirm your forum-authentication account\r\n\r\n"
	body := fmt.Sprintf("Welcome!\nClick the link to verify your account:\n%s", link)
	msg := []byte(subject + body)

	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")

	// Set custom dialer with timeout
	client, err := net.DialTimeout("tcp", "smtp.gmail.com:587", 10*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}
	defer client.Close()

	return smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, msg)
}
