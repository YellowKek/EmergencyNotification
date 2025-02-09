package auth

import (
	"app/internal/model"
	"app/internal/service"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
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
	regReq := RegisterRequest{}
	if err := c.BodyParser(&regReq); err != nil {
		return fmt.Errorf("body parser: %w", err)
	}
	exists, _ := h.Storage.UserService.GetByEmail(regReq.Email)

	if exists != nil {
		return errors.New("the user already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(regReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("bcrypt hash: %w", err)
	}

	user := model.User{
		Email:    regReq.Email,
		Name:     regReq.Name,
		Surname:  regReq.Surname,
		Password: string(hashedPassword),
	}

	err = h.Storage.UserService.CreateUser(user)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusCreated)
}
