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
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
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
	if client == nil {
		log.Fatal("Ошибка получения клиента OAuth 2.0")
	}

	// Создание нового сервиса Gmail
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Ошибка создания сервиса Gmail: %v", err)
	}

	// Формирование письма
	to := "galimardanova123@gmail.com"
	subject := "Subject: Test Email"
	body := "This is a test email sent from Go!"
	message := []byte(fmt.Sprintf("To: %s\r\n%s\r\n\r\n%s", to, subject, body))

	// Правильное кодирование в Base64
	encodedMessage := base64.StdEncoding.EncodeToString(message)

	// Отправка письма
	msg := &gmail.Message{
		Raw: encodedMessage,
	}
	_, err = srv.Users.Messages.Send("me", msg).Do()
	if err != nil {
		log.Printf("Ошибка отправки письма: %v", err)
		return
	}

	fmt.Println("Email отправлен успешно!")
}

// getClient получает токен доступа, используя OAuth 2.0
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	tok, err := tokenFromFile("token.json")
	if err != nil {
		log.Println("Токен не найден, требуется авторизация.")
		tok = getTokenFromWeb(config)
		err := saveToken("token.json", tok)
		if err != nil {
			log.Printf("Ошибка сохранения токена: %v", err)
			return nil
		}
	}
	return config.Client(ctx, tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	// Создаем канал для передачи кода авторизации
	codeChan := make(chan string)

	// HTTP-сервер для обработки редиректа
	http.HandleFunc("/oauth2callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Код авторизации не найден", http.StatusBadRequest)
			return
		}
		fmt.Fprintln(w, "Авторизация успешна! Теперь вернитесь в консоль.")
		codeChan <- code
	})

	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil)) // Запускаем сервер
	}()

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Перейдите по ссылке для авторизации: %s\n", authURL)

	code := <-codeChan // Ждем код из канала

	tok, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("Ошибка обмена кода на токен: %v", err)
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
