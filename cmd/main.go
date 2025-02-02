package main

import (
	"EmergencyNotification/internal/config"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

func main() {
	cfgYaml, err := config.ParseFromYaml()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Создание конфигурации OAuth 2.0
	cfg := &oauth2.Config{
		ClientID:     cfgYaml.SMTP.ClientID,
		ClientSecret: cfgYaml.SMTP.ClientSecret,
		RedirectURL:  "http://localhost:8080/oauth2callback",
		Scopes:       []string{gmail.GmailSendScope},
		Endpoint:     google.Endpoint,
	}

	// Получение токена доступа
	client := getClient(ctx, cfg)

	// Создание нового сервиса Gmail
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	// Создание сообщения
	to := "damirgarifullin7@gmail.com"
	subject := "Subject: Test Email"
	body := "This is a test email sent from Go!"
	message := []byte(fmt.Sprintf("To: %s\r\n%s\r\n\r\n%s", to, subject, body))

	// Кодирование сообщения в base64
	encodedMessage := base64.URLEncoding.EncodeToString(message)

	// Отправка письма
	msg := &gmail.Message{
		Raw: encodedMessage,
	}
	_, err = srv.Users.Messages.Send("me", msg).Do()
	if err != nil {
		log.Fatalf("Unable to send email: %v", err)
	}

	fmt.Println("Email sent successfully!")
}

// getClient получает токен доступа, используя OAuth 2.0
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	// Загружаем токен из файла, если он существует
	tok, err := tokenFromFile("token.json")
	if err != nil {
		tok = getTokenFromWeb(config)
		err := saveToken("token.json", tok)
		if err != nil {
			log.Print("Unable to save token:", err)
			return nil
		}
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb получает токен доступа от пользователя
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	// Генерируем URL для авторизации
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Перейдите по следующему URL для авторизации: \n%v\n", authURL)

	// Получаем код авторизации от пользователя
	var code string
	fmt.Print("Введите код авторизации: ")
	fmt.Scan(&code)

	// Обмениваем код на токен
	tok, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// saveToken сохраняет токен в файл
func saveToken(filename string, token *oauth2.Token) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	// Кодируем токен в JSON и записываем в файл
	return json.NewEncoder(f).Encode(token)
}

// tokenFromFile загружает токен из файла
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var tok oauth2.Token
	if err := json.NewDecoder(f).Decode(&tok); err != nil {
		return nil, err
	}
	return &tok, nil
}
