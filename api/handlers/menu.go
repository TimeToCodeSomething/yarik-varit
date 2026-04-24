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

// GetMenuItems возвращает все позиции меню, опционально с фильтром по категории
func (h *Handler) GetMenuItems(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")

	var rows *sql.Rows
	var err error

	if category != "" {
		rows, err = h.db.Query("SELECT id, name, price, vol, category FROM menu WHERE category = $1 ORDER BY id", category)
	} else {
		rows, err = h.db.Query("SELECT id, name, price, vol, category FROM menu ORDER BY category, id")
	}

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

// GetMenuByCategories возвращает меню сгруппированное по категориям
func (h *Handler) GetMenuByCategories(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query("SELECT id, name, price, vol, category FROM menu ORDER BY category, id")
	if err != nil {
		http.Error(w, "Ошибка выполнения запроса", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Группируем по категориям
	menuByCategory := make(map[string][]models.MenuItem)
	categories := []string{}
	categorySet := make(map[string]bool)

	for rows.Next() {
		var item models.MenuItem
		if err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.Vol, &item.Category); err != nil {
			http.Error(w, "Ошибка чтения данных", http.StatusInternalServerError)
			return
		}

		if !categorySet[item.Category] {
			categories = append(categories, item.Category)
			categorySet[item.Category] = true
		}

		menuByCategory[item.Category] = append(menuByCategory[item.Category], item)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Ошибка итерации", http.StatusInternalServerError)
		return
	}

	// Возвращаем структурированный ответ
	response := make(map[string][]models.MenuItem)
	for _, cat := range categories {
		response[cat] = menuByCategory[cat]
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

// CreateMenuItem добавляет новую позицию в меню
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

// UpdateMenuItem обновляет позицию меню
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

	affected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Ошибка проверки обновления", http.StatusInternalServerError)
		return
	}
	if affected == 0 {
		http.Error(w, "Позиция не найдена", http.StatusNotFound)
		return
	}

	item.ID = idInt
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(item)
}

// DeleteMenuItem удаляет позицию из меню
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

	affected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Ошибка проверки удаления", http.StatusInternalServerError)
		return
	}
	if affected == 0 {
		http.Error(w, "Позиция не найдена", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
