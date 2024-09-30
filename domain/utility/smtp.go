// Reference: https://www.loginradius.com/blog/engineering/sending-emails-with-golang/
package utility

import (
	"crypto/tls"
	"fmt"
	"os"

	gomail "gopkg.in/mail.v2"
)

// func main() {
// 	godotenv_err := godotenv.Load()
// 	if godotenv_err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	// Sender data.
// 	sender := os.Getenv("EMAIL_SENDER")
// 	password := os.Getenv("EMAIL_PASSWORD")

// 	// Receiver email address.
// 	receiver := []string{
// 		os.Getenv("RECEIVER"),
// 	}

// 	// Message.
// 	message := []byte("Hello,This is a simple mail. There is only text, no attachments are there. The mail is sent using Golang. Thank You!!!")
// 	// send_email_smtp(sender, password, receiver, message)

// 	send_email_gomail(sender, password, receiver, message)

// }

func SendEmail(receiver []string, subject string, message string) error {
	if len(receiver) == 0 || subject == "" {
		return fmt.Errorf("receiver or subject is empty")
	}

	// Sender data.
	sender := os.Getenv("EMAIL_ADDRESS")
	password := os.Getenv("EMAIL_PASSWORD")
	// smtp server configuration.
	// 	smtpHost := "smtp.spaceai.jp"
	// 	// smtpHost := "smtp.gmail.com"
	// 	smtpPort := "587"

	// 	// Authentication.
	// 	auth := smtp.PlainAuth("", sender, password, smtpHost)

	// 	// Sending email.
	// 	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, sender, receiver, message)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	// 	fmt.Println("Email Sent Successfully!")
	// }

	mail := gomail.NewMessage()

	// Set E-Mail sender
	mail.SetHeader("From", sender)

	// Set E-Mail receivers
	mail.SetHeader("To", receiver[0])

	// Set E-Mail subject
	mail.SetHeader("Subject", subject)

	// Set E-Mail body. You can set plain text or html with text/html
	mail.SetBody("text/plain", message)

	// Settings for SMTP server
	server := gomail.NewDialer("smtp.gmail.com", 587, sender, password)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := server.DialAndSend(mail); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
