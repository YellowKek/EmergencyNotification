package main

import (
	"app/internal"
	"app/internal/auth"
	"app/internal/repository"
	"app/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

var kafkaBrokers = []string{"localhost:9092"}

//const KafkaTopic =

func main() {
	logrus.Print("app start")
	db, err := internal.NewDB()
	if err != nil {
		logrus.Fatal(err)
	}

	app := fiber.New()

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)

	err = userService.CreateTables()
	if err != nil {
		logrus.Fatal(err)
	}

	authHandler := &auth.AuthHandler{Storage: &auth.AuthStorage{UserService: userService}}

	app.Post("/register", authHandler.Register)

	logrus.Fatal(app.Listen(":8080"))
	if err != nil {
		logrus.Fatal(err)
	}
	defer db.Close()
}
