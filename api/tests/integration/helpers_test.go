package integration_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"yarik-varit/api/handlers"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

var (
	testDB     *sql.DB
	testServer *httptest.Server
	adminToken string
)

func TestMain(m *testing.M) {
	// Загружаем .env из корня проекта
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Println("Файл .env не найден, используются переменные окружения")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("Переменная DB_URL не задана")
	}

	var err error
	testDB, err = sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	if err := testDB.Ping(); err != nil {
		log.Fatal("Ошибка пинга БД:", err)
	}

	// Поднимаем тестовый сервер с реальным роутером
	testServer = httptest.NewServer(handlers.NewRouter(testDB))

	// Получаем токен админа через реальный /login
	adminUser := os.Getenv("TEST_ADMIN_USER")
	adminPass := os.Getenv("TEST_ADMIN_PASS")
	if adminUser == "" || adminPass == "" {
		log.Fatal("Переменные TEST_ADMIN_USER и TEST_ADMIN_PASS не заданы")
	}
	adminToken = mustLogin(adminUser, adminPass)

	code := m.Run()

	testServer.Close()
	testDB.Close()
	os.Exit(code)
}

// mustLogin делает запрос на /login и возвращает токен
func mustLogin(username, password string) string {
	body, _ := json.Marshal(map[string]string{"username": username, "password": password})
	res, err := http.Post(testServer.URL+"/login", "application/json", bytes.NewBuffer(body))
	if err != nil || res.StatusCode != http.StatusOK {
		log.Fatalf("Не удалось залогиниться как %s: %v", username, err)
	}
	defer res.Body.Close()
	var data map[string]string
	json.NewDecoder(res.Body).Decode(&data)
	return data["token"]
}

// get выполняет GET запрос с опциональным токеном
func get(path, token string) *http.Response {
	req, _ := http.NewRequest(http.MethodGet, testServer.URL+path, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("Ошибка запроса:", err)
	}
	return res
}

// post выполняет POST запрос с опциональным токеном
func post(path, token string, body interface{}) *http.Response {
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, testServer.URL+path, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("Ошибка запроса:", err)
	}
	return res
}

// decode декодирует JSON тело ответа
func decode(res *http.Response, v interface{}) {
	defer res.Body.Close()
	json.NewDecoder(res.Body).Decode(v)
}

// bodyString читает тело ответа как строку
func bodyString(res *http.Response) string {
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	return string(b)
}

// doRequest выполняет запрос с произвольным методом
func doRequest(method, path, token string, body interface{}) (*http.Response, error) {
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, testServer.URL+path, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return http.DefaultClient.Do(req)
}
