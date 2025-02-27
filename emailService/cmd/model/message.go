package model

type KafkaMessage struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Location string `json:"location"`
	Email    string `json:"email"`
}
