package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"yarik-varit/api/handlers"
	"yarik-varit/api/models"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, используются переменные окружения")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("Переменная DB_URL не задана")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer db.Close()

	// Проверка подключения
	if err := db.Ping(); err != nil {
		log.Fatal("Ошибка пинга БД:", err)
	}

	log.Println("✅ Подключение к БД успешно")

	h := handlers.NewHandler(db)

	// Используем свой мультиплексор для маршрутов
	mux := http.NewServeMux()

	// Публичные эндпоинты
	mux.HandleFunc("POST /login", h.Login)
	mux.HandleFunc("GET /menu", h.GetMenuItems)
	mux.HandleFunc("GET /menu/categories", h.GetMenuByCategories)

	// Ограничение ролей для безопасности
	adminOnly := []models.Role{models.RoleAdmin}
	staffRoles := []models.Role{models.RoleAdmin, models.RoleBarista}
	anyUser := []models.Role{models.RoleAdmin, models.RoleBarista, models.RoleClient}

	// Menu endpoints
	mux.HandleFunc("POST /menu", handlers.AuthMiddleware(adminOnly, h.CreateMenuItem))
	mux.HandleFunc("PUT /menu/{id}", handlers.AuthMiddleware(adminOnly, h.UpdateMenuItem))
	mux.HandleFunc("DELETE /menu/{id}", handlers.AuthMiddleware(adminOnly, h.DeleteMenuItem))

	// Orders endpoints
	mux.HandleFunc("GET /orders", handlers.AuthMiddleware(staffRoles, h.GetOrders))
	mux.HandleFunc("GET /orders/{id}", handlers.AuthMiddleware(anyUser, h.GetOrderByID))
	mux.HandleFunc("POST /orders", handlers.AuthMiddleware(anyUser, h.CreateOrder))
	mux.HandleFunc("PATCH /orders/{id}/status", handlers.AuthMiddleware(staffRoles, h.UpdateOrderStatus))
	mux.HandleFunc("DELETE /orders/{id}", handlers.AuthMiddleware(staffRoles, h.DeleteOrder))

	log.Printf("🚀 Сервер 'Ярик Варит' запущен на :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
