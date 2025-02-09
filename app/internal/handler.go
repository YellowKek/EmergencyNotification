package internal

import (
	"app/internal/service"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/gofiber/fiber/v2"
	"os"
	"os/signal"
	"strings"
)

type EmergencyHandler struct {
	s *service.UserService
	p sarama.AsyncProducer
}

func NewEmergencyHandler(s *service.UserService, p sarama.AsyncProducer) *EmergencyHandler {
	return &EmergencyHandler{s: s, p: p}
}

func (h *EmergencyHandler) EmergencyCall(c *fiber.Ctx) error {
	var body map[string]int32
	if err := json.Unmarshal(c.Body(), &body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"body parser": err.Error()})
	}

	userID, exists := body["user_id"]
	if !exists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"body user_id": nil})
	}
	
	emergencyGroup, err := h.s.GetEmergencyGroups(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"getting emergency groups": err.Error()})
	}

	user, err := h.s.GetByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"getting user": err.Error()})
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for group, value := range emergencyGroup {
		if strings.Compare(group, "email") == 0 {
			kafkaMessage, err := json.Marshal(struct {
				Sender   string `json:"sender"`
				Receiver string `json:"receiver"`
			}{
				Sender:   user.Name + user.Surname,
				Receiver: value,
			})
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"marshalling error": err.Error()})
			}

			ProduceMessage(EmailTopic, h.p, kafkaMessage, signals)
		}
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "ok"})
}

type addGroupDTO struct {
	UserID int32  `json:"user_id"`
	Group  string `json:"group"`
	Value  string `json:"value"`
}

func (h *EmergencyHandler) AddEmergencyGroup(c *fiber.Ctx) error {
	var req addGroupDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error parsing user": err.Error()})
	}
	user, err := h.s.GetByID(req.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error getting user": err.Error()})
	}

	err = h.s.AddEmergencyGroup(user.Id, req.Group, req.Value)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error adding group": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}
