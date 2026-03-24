package main

import (
	"fmt"
	"net/http"
	"yarik-varit/api/handlers"
)

func main() {
	fmt.Println("Ярик Варит! Сервер запустился!")
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/menu", handlers.MenuHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Добро пожаловать в Ярик Варит!")
}
