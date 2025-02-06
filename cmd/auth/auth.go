package main

import (
	"auth/database/models"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func main() {
	app := fiber.New()

	authHandler := &AuthHandler{&AuthStorage{map[string]models.User{}}}

	app.Post("/register", authHandler.Register)
	app.Post("/login", authHandler.Login)

	logrus.Fatal(app.Listen(":8080"))
}
