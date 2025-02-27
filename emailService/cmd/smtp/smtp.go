package smtp

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gmailService/model"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

var (
	cfgYaml = getCfgFromYaml()

	// Создание конфигурации OAuth 2.0
	config = &oauth2.Config{
		ClientID:     cfgYaml.SMTP.ClientID,
		ClientSecret: cfgYaml.SMTP.ClientSecret,
		RedirectURL:  "http://localhost:8080/oauth2callback",
		Scopes:       []string{gmail.GmailSendScope},
		Endpoint:     google.Endpoint,
	}
)

func getCfgFromYaml() *Config {
	cfgYaml, err := ParseFromYaml()
	if err != nil {
		logrus.Printf("Ошибка загрузки конфигурации: %v", err)
	}
	return cfgYaml
}

func GetGmailService() (*gmail.Service, error) {
	ctx := context.Background()
	client, err := getClient(ctx, config)
	if err != nil {
		return nil, err
	}
	if client == nil {
		logrus.Print("Ошибка получения клиента OAuth 2.0")
	}
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		logrus.Printf("Ошибка создания сервиса Gmail: %v", err)
		return nil, err
	}
	return srv, nil
}

// SendEmail TODO
// исправить так чтобы вместо почты подставлялось имя при регистрации
func SendEmail(srv *gmail.Service, message model.KafkaMessage) error {
	// Формирование письма
	to := message.Receiver
	subject := fmt.Sprintf("ATTENTION!!! %s is in danger", message.Sender)
	body := fmt.Sprintf("Contact the person immediately (%s) and call 911\nThe address from which the message was sent: %s", message.Email, message.Location)

	// Создание MIME-сообщения
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Заголовки письма
	headers := textproto.MIMEHeader{}
	headers.Set("To", to)
	headers.Set("Subject", "=?UTF-8?B?"+base64.StdEncoding.EncodeToString([]byte(subject))+"?=") // Кодируем заголовок
	headers.Set("MIME-Version", "1.0")
	headers.Set("Content-Type", fmt.Sprintf("multipart/alternative; boundary=%s", writer.Boundary()))

	// Запись заголовков
	_, err := writer.CreatePart(headers)
	if err != nil {
		logrus.Errorf("Ошибка создания заголовков: %v", err)
		return err
	}

	// Тело письма (текстовая часть)
	part, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type": []string{"text/plain; charset=UTF-8"},
	})
	if err != nil {
		logrus.Errorf("Ошибка создания текстовой части: %v", err)
		return err
	}
	_, err = part.Write([]byte(body))
	if err != nil {
		logrus.Errorf("Ошибка записи текстовой части: %v", err)
		return err
	}

	// Завершение формирования MIME-сообщения
	err = writer.Close()
	if err != nil {
		logrus.Errorf("Ошибка завершения MIME-сообщения: %v", err)
		return err
	}

	// Кодирование сообщения в Base64
	encodedMessage := base64.URLEncoding.EncodeToString(buf.Bytes())

	// Отправка письма
	msg := &gmail.Message{
		Raw: strings.TrimRight(encodedMessage, "="), // Убираем лишние символы '='
	}
	_, err = srv.Users.Messages.Send("me", msg).Do()
	if err != nil {
		logrus.Errorf("Ошибка отправки письма: %v", err)
		return err
	}

	logrus.Info("Email отправлен успешно!")
	return nil
}

// getClient получает токен доступа, используя OAuth 2.0
func getClient(ctx context.Context, config *oauth2.Config) (*http.Client, error) {
	tok, err := tokenFromFile("token.json")
	if err != nil {
		logrus.Print("Токен не найден, требуется авторизация.")
		tok = getTokenFromWeb()
		err := saveToken("token.json", tok)
		if err != nil {
			logrus.Printf("Ошибка сохранения токена: %v", err)
			return nil, err
		}
	}
	return config.Client(ctx, tok), nil
}

func getTokenFromWeb() *oauth2.Token {
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
