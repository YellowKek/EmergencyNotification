package main

import (
	"github.com/sirupsen/logrus"
	"gmailService/smtp"
	"net/http"
)

func main() {
	service, err := smtp.GetGmailService()
	if err != nil {
		logrus.Error(err)
	}

	err = smtp.SendEmail(service, "damirgarifullin7@gmail.com", "damirgarifullin7@gmail.com")
	if err != nil {
		return
	}
	err = http.ListenAndServe(":8081", nil)
	if err != nil {
		panic(err)
	}
}
