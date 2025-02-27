package auth

import (
	"app/internal"
	"app/internal/model"
	"app/internal/service"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type (
	AuthHandler struct {
		Storage *AuthStorage
	}
	AuthStorage struct {
		UserService *service.UserService
	}
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	logrus.Print("[AuthHandler] Register request")
	regReq := RegisterRequest{}
	if err := c.BodyParser(&regReq); err != nil {
		logrus.Print(err.Error())
		return fmt.Errorf("body parser: %w", err)
	}
	exists, _ := h.Storage.UserService.GetByEmail(regReq.Email)

	if exists != nil {
		return errors.New("the user already exists")
	}

	user := model.User{
		Email:    regReq.Email,
		Name:     regReq.Name,
		Surname:  regReq.Surname,
		Password: regReq.Password,
	}

	if !internal.ValidateUser(user) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid data"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(regReq.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Print(err.Error())
		return fmt.Errorf("bcrypt hash: %w", err)
	}

	user.Password = string(hashedPassword)

	user, err = h.Storage.UserService.CreateUser(user)
	if err != nil {
		logrus.Print(err.Error())
		return err
	}

	return c.Status(fiber.StatusOK).JSON(user)
}
