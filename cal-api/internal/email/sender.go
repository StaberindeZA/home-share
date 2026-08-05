package email

import (
	"log"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type GomailSender struct {
	messages GomailMessages
}

func (es GomailSender) getDialer() *gomail.Dialer {
	host := os.Getenv("EMAIL_HOST")
	portString := os.Getenv("EMAIL_PORT")
	user := os.Getenv("EMAIL_USER")
	pass := os.Getenv("EMAIL_PASS")

	port, err := strconv.Atoi(portString)
	if err != nil {
		log.Fatal("Error converting port to int")
	}

	return gomail.NewDialer(host, port, user, pass)
}

func (es GomailSender) SendOtp(to, otp string) error {
	message, err := es.messages.SendOtp(to, otp)
	if err != nil {
		log.Printf("Failed to retrieve SendOtp message: %v", err)
		return err
	}

	dialer := es.getDialer()
	if err := dialer.DialAndSend(message); err != nil {
		log.Printf("Failed to send SendOtp mail: %v", err)
		return err
	} else {
		log.Printf("Successfully sent SendOtp email to: %s", to)
	}

	return nil
}

func NewGomailSender(messages GomailMessages) GomailSender {
	return GomailSender{
		messages: messages,
	}
}
