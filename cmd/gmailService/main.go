package main

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	service, err := getGmailService()
	if err != nil {
		logrus.Error(err)
	}

	err = sendEmail(service, "damirgarifullin7@gmail.com", "damirgarifullin7@gmail.com")
	if err != nil {
		return
	}
	err = http.ListenAndServe(":8081", nil)
	if err != nil {
		panic(err)
	}
}
