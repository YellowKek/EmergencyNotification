package auth

//
//import (
//	"errors"
//	"fmt"
//	"github.com/gofiber/fiber/v2"
//	"github.com/golang-jwt/jwt/v5"
//	"github.com/sirupsen/logrus"
//)
//
//type LoginRequest struct {
//	Email    string `json:"email"`
//	Password string `json:"password"`
//}
//
//type LoginResponse struct {
//	AccessToken string `json:"access_token"`
//}
//
//var (
//	errBadCredentials = errors.New("email or password is incorrect")
//)
//
//var jwtSecretKey = []byte("secret-key")
//
//func (h *AuthHandler) Login(c *fiber.Ctx) error {
//	regReq := LoginRequest{}
//	if err := c.BodyParser(&regReq); err != nil {
//		return fmt.Errorf("body parser: %w", err)
//	}
//
//	user, exists := h.Storage.users[regReq.Email]
//	if !exists {
//		return errBadCredentials
//	}
//	if user.GetPassword() != regReq.Password {
//		return errBadCredentials
//	}
//
//	payload := jwt.MapClaims{
//		"sub": user.Email,
//		//"exp": time.Now().Add(time.Hour * 72).Unix(),
//	}
//
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
//
//	t, err := token.SignedString(jwtSecretKey)
//	if err != nil {
//		logrus.WithError(err).Error("JWT token signing")
//		return c.SendStatus(fiber.StatusInternalServerError)
//	}
//	logrus.Print("successfully logged in")
//
//	return c.JSON(LoginResponse{AccessToken: t})
//}
