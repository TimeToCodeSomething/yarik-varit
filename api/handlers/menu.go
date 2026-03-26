package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"yarik-varit/api/models"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func MenuHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("pgx", "postgres://kuimovmihail@localhost:5432/yarik_varit?sslmode=disable")
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}

	rows, err := db.Query("SELECT id, name, price, vol, category FROM menu")
	if err != nil {
		http.Error(w, "Ошибка выполнения запроса", http.StatusInternalServerError)
		return
	}

	var items []models.MenuItem
	for rows.Next() {
		var item models.MenuItem
		rows.Scan(&item.ID, &item.Name, &item.Price, &item.Vol, &item.Category)
		items = append(items, item)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(items)
}
