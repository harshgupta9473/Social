package utils

import (
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func SendVerificationEmail(email,  token string) error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	from := os.Getenv("emailID")
	password := os.Getenv("apppassword")
	smtpHost := os.Getenv("smtpHost")
	smtpPort := os.Getenv("smtpPort")

	link := os.Getenv("link") + "/verify?token=" + token + "&userID=" + email 

	msg := "From: " + from + "\n" +
		"To: " + email + "\n" +
		"Subject: Email Verification \n\n" +
		"Click the link to verify your email: " + link

	auth := smtp.PlainAuth("", from, password, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{email}, []byte(msg))
}
