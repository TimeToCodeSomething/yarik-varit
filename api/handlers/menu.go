package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"yarik-varit/api/models"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) GetMenuItems(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query("SELECT id, name, price, vol, category FROM menu")
	if err != nil {
		http.Error(w, "Ошибка выполнения запроса", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []models.MenuItem
	for rows.Next() {
		var item models.MenuItem
		if err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.Vol, &item.Category); err != nil {
			http.Error(w, "Ошибка чтения данных", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Ошибка итерации", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(items)
}

func (h *Handler) CreateMenuItem(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Ошибка вставки", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func (h *Handler) UpdateMenuItem(w http.ResponseWriter, r *http.Request) {
	idInt, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Неверный id", http.StatusBadRequest)
		return
	}

	var item models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Ошибка чтения запроса", http.StatusBadRequest)
		return
	}

	result, err := h.db.Exec(
		"UPDATE menu SET name=$1, price=$2, vol=$3, category=$4 WHERE id=$5",
		item.Name, item.Price, item.Vol, item.Category, idInt,
	)
	if err != nil {
		http.Error(w, "Ошибка обновления", http.StatusInternalServerError)
		return
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "Позиция не найдена", http.StatusNotFound)
		return
	}

	item.ID = idInt
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(item)
}

func (h *Handler) DeleteMenuItem(w http.ResponseWriter, r *http.Request) {
	idInt, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Неверный id", http.StatusBadRequest)
		return
	}

	result, err := h.db.Exec("DELETE FROM menu WHERE id=$1", idInt)
	if err != nil {
		http.Error(w, "Ошибка удаления", http.StatusInternalServerError)
		return
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "Позиция не найдена", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
