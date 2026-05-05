package main

import (
	"database/sql"
	"log"
	"net/http"
	"yarik-varit/api/handlers"
	"yarik-varit/api/models"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// Подключение к БД
	db, err := sql.Open("pgx", "postgres://kuimovmihail:password@localhost:5432/yarik_varit")
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

	log.Println("🚀 Сервер 'Ярик Варит' запущен на :8080")

	// ВАЖНО: передаем созданный mux вместо nil
	log.Fatal(http.ListenAndServe(":8080", mux))
}
