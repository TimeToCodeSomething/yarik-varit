package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"yarik-varit/api/handlers"
)

func main() {
	fmt.Println("Ярик Варит! Сервер запустился!")
	db, err := sql.Open("pgx", "postgres://kuimovmihail@localhost:5432/yarik_varit?sslmode=disable")
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	h := handlers.NewHandler(db)

	http.HandleFunc("GET /", homeHandler)
	http.HandleFunc("GET /menu", h.GetMenuItems)
	http.HandleFunc("POST /menu", h.CreateMenuItem)
	http.HandleFunc("PUT /menu/{id}", h.UpdateMenuItem)
	http.HandleFunc("DELETE /menu/{id}", h.DeleteMenuItem)

	http.HandleFunc("GET /orders", h.GetOrders)
	http.HandleFunc("POST /orders", h.CreateOrder)
	http.HandleFunc("GET /orders/{id}", h.GetOrderByID)
	http.HandleFunc("PATCH /orders/{id}/status", h.UpdateOrderStatus)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Добро пожаловать в Ярик Варит!")
}
