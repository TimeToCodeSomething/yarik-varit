package main

import (
	"database/sql"
	"log"
	"net/http"
	"yarik-varit/api/handlers"

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

	// Menu endpoints
	http.HandleFunc("GET /menu", h.GetMenuItems)
	http.HandleFunc("GET /menu/categories", h.GetMenuByCategories)
	http.HandleFunc("POST /menu", h.CreateMenuItem)
	http.HandleFunc("PUT /menu/{id}", h.UpdateMenuItem)
	http.HandleFunc("DELETE /menu/{id}", h.DeleteMenuItem)

	// Orders endpoints
	http.HandleFunc("GET /orders", h.GetOrders)
	http.HandleFunc("GET /orders/{id}", h.GetOrderByID)
	http.HandleFunc("POST /orders", h.CreateOrder)
	http.HandleFunc("PATCH /orders/{id}/status", h.UpdateOrderStatus)

	log.Println("🚀 Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
