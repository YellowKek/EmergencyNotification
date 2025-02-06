package main

import (
	"auth/database/models"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

var id int32 = 1

type (
	AuthHandler struct {
		storage *AuthStorage
	}
	AuthStorage struct {
		users map[string]models.User
	}
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	regReq := RegisterRequest{}
	if err := c.BodyParser(&regReq); err != nil {
		return fmt.Errorf("body parser: %w", err)
	}

	if _, exists := h.storage.users[regReq.Email]; exists {
		return errors.New("the user already exists")
	}

	user := models.User{
		Id:      id,
		Email:   regReq.Email,
		Name:    regReq.Name,
		Surname: regReq.Surname,
	}

	user.SetPassword(regReq.Password)
	h.storage.users[regReq.Email] = user
	id++

	logrus.Printf("%+v", h.storage.users)

	return c.SendStatus(fiber.StatusCreated)
}
