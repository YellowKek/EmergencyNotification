package main

import (
	"app/internal"
	"app/internal/auth"
	"app/internal/repository"
	"app/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Print("app start")
	db, err := internal.NewDB()
	if err != nil {
		logrus.Fatal(err)
	}

	producer, err := internal.SetupProducer()
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Print("producer start")

	app := fiber.New()

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)

	err = userService.CreateTables()
	if err != nil {
		logrus.Fatal(err)
	}

	authHandler := &auth.AuthHandler{Storage: &auth.AuthStorage{UserService: userService}}
	emergencyHandler := internal.NewEmergencyHandler(userService, producer)

	app.Post("/register", authHandler.Register)
	app.Post("/emergency", emergencyHandler.EmergencyCall)
	app.Post("/addGroup", emergencyHandler.AddEmergencyGroup)

	logrus.Fatal(app.Listen(":8080"))
	if err != nil {
		logrus.Fatal(err)
	}
	defer db.Close()
}
