package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"yarik-varit/api/models"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Handler struct {
	db *sql.DB // Приватное поле, так как с маленькой буквы
}

func NewHandler(db *sql.DB) *Handler { // Создаем конструктор для структуры Handler
	return &Handler{db: db}
}

func (h *Handler) MenuHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rows, err := h.db.Query("SELECT id, name, price, vol, category FROM menu")
		if err != nil {
			http.Error(w, "Ошибка выполнения запроса", http.StatusInternalServerError)
			return
		}
		defer rows.Close() // Закрываем rows после использования
		// (defer закрывает после того, как функция прекратит работу)

		var items []models.MenuItem
		for rows.Next() {
			var item models.MenuItem
			if err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.Vol, &item.Category); err != nil {
				http.Error(w, "Ошибка чтения данных", http.StatusInternalServerError)
				return
			}
			items = append(items, item)
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(items)
	case http.MethodPost:
		var item models.MenuItem
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			http.Error(w, "Ошибка чтения запроса", http.StatusBadRequest)
			return
		}
		err := h.db.QueryRow(
			"INSERT INTO menu (name, price, vol, category) VALUES ($1, $2, $3, $4) RETURNING id",
			item.Name, item.Price, item.Vol, item.Category,
		).Scan(&item.ID)
		if err != nil {
			http.Error(w, "Ошибка вставки", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(item)

	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) MenuItemHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/menu/")
	switch r.Method {
	case http.MethodPut:
		var item models.MenuItem
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			http.Error(w, "Ошибка изменения", http.StatusBadRequest)
			return
		}
		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Неверный id", http.StatusBadRequest)
			return
		}

		_, err = h.db.Exec(
			"UPDATE menu SET name=$1, price=$2, vol=$3, category=$4 WHERE id=$5",
    		item.Name, item.Price, item.Vol, item.Category, idInt,
		)
	case http.MethodDelete:

	default:
	}
}
