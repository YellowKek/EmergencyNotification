package main

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"gmailService/model"
	"gmailService/smtp"
	"log"
	"sync"
)

const KafkaTopic = "email"

func main() {
	consumer, err := sarama.NewConsumer([]string{"kafka:29091"}, nil)
	if err != nil {
		logrus.Fatalf("Error creating Sarama consumer: %v", err)
	}
	defer consumer.Close()

	partConsumer, err := consumer.ConsumePartition(KafkaTopic, 0, 0)
	if err != nil {
		log.Fatalf("Failed to start partition consumer: %v", err)
	}
	defer partConsumer.Close()

	messageService := NewMessageService()

	go ConsumeMessage(messageService, partConsumer)
	select {}
}

func ConsumeMessage(messageService *MessageService, partConsumer sarama.PartitionConsumer) {
	for {
		log.Print("start read msg")
		select {
		case msg, ok := <-partConsumer.Messages():
			if !ok {
				log.Print("Channel closed")
				return
			}
			messageService.ProcessMessage(msg)
		}
	}
}

type MessageService struct {
	mu           sync.Mutex
	partConsumer sarama.PartitionConsumer
}

func NewMessageService() *MessageService {
	return &MessageService{}
}

func (s *MessageService) ProcessMessage(msg *sarama.ConsumerMessage) {
	s.mu.Lock()

	var message model.KafkaMessage
	err := json.Unmarshal(msg.Value, &message)
	if err != nil {
		logrus.Print(err)
		return
	}
	log.Print("Message consumed: ", string(msg.Value))

	service, err := smtp.GetGmailService()
	if err != nil {
		logrus.Error(err)
	}

	err = smtp.SendEmail(service, message)
	if err != nil {
		return
	}
	logrus.Print("email sent")
	s.mu.Unlock()
}
