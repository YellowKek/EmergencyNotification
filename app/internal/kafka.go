package internal

import (
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"os"
)

var kafkaBrokers = []string{"kafka:29091"}

const EmailTopic = "email"

func SetupProducer() (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	return sarama.NewAsyncProducer(kafkaBrokers, config)
}

func ProduceMessage(kafkaTopic string, producer sarama.AsyncProducer, message []byte, signals chan os.Signal) {
	logrus.Print("sending message to kafka")
	kafkaMessage := &sarama.ProducerMessage{
		Topic: kafkaTopic,
		Value: sarama.ByteEncoder(message),
	}
	select {
	case producer.Input() <- kafkaMessage:
		logrus.Print("message produced")
	case signal := <-signals:
		producer.AsyncClose()
		logrus.Print("close signal", signal)
		return
	}
}
