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
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/menu", h.MenuHandler)
	http.HandleFunc("/menu/", h.MenuItemHandler)
	http.HandleFunc("/orders", h.OrderHandler)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Добро пожаловать в Ярик Варит!")
}
