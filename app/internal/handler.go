package internal

import (
	"app/internal/service"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"strconv"
	"strings"
)

type EmergencyHandler struct {
	s *service.UserService
	p sarama.AsyncProducer
}

func NewEmergencyHandler(s *service.UserService, p sarama.AsyncProducer) *EmergencyHandler {
	return &EmergencyHandler{s: s, p: p}
}

type emergencyCallDTO struct {
	UserID   int32  `json:"user_id"`
	Location string `json:"location"`
}

func (h *EmergencyHandler) EmergencyCall(c *fiber.Ctx) error {
	logrus.Print("emergency call")
	var body emergencyCallDTO
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"body parser": err.Error()})
	}
	logrus.Print("user id:", body.UserID)

	emergencyGroup, err := h.s.GetEmergencyGroups(body.UserID)
	if err != nil {
		logrus.Print(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"getting emergency groups": err.Error()})
	}

	user, err := h.s.GetByID(body.UserID)
	if err != nil {
		logrus.Print(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"getting user": err.Error()})
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for group, value := range emergencyGroup {
		if strings.Compare(group, "email") == 0 {
			kafkaMessage, err := json.Marshal(struct {
				Sender   string `json:"sender"`
				Receiver string `json:"receiver"`
				Location string `json:"location"`
				Email    string `json:"email"`
			}{
				Sender:   user.Name + " " + user.Surname,
				Receiver: value,
				Location: body.Location,
				Email:    user.Email,
			})
			logrus.Printf("message to receive: %s %s %s %s %s", user.Name, user.Surname, value, body.Location, user.Email)
			if err != nil {
				logrus.Print(err.Error())
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error parsing": err.Error()})
	}

	if req.Group == "email" {
		if !ValidateEmail(req.Value) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "error email is invalid"})
		}
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

func (h *EmergencyHandler) GetEmergencyGroups(c *fiber.Ctx) error {
	logrus.Print("get groups request")

	userIDString := c.Query("user_id")
	if userIDString == "" {
		logrus.Print("user id not found in request")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id is required"})
	}

	logrus.Print("user id:", userIDString)

	userId, err := strconv.ParseInt(userIDString, 10, 32)
	if err != nil {
		logrus.Print("failed to convert user id to int")
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	_, err = h.s.GetByID(int32(userId))

	if err != nil {
		logrus.Print("user id not found")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "user not found"})
	}

	groups, err := h.s.GetEmergencyGroups(int32(userId))
	if err != nil {
		logrus.Print("failed to get groups " + err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(groups)
}
