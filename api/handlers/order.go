package handlers

import (
	"net/http"
	"yarik-varit/api/models"
	"encoding/json"
)

func (h *Handler) OrderHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:	// GET
		rows, err := h.db.Query("SELECT id, tm, status FROM orders")
		if err != nil {
			http.Error(w, "Ошибка выполнения запроса", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var items []models.Order
		for rows.Next() {
			var item models.Order
			if err := rows.Scan(&item.ID, &item.Time, &item.Status); err != nil {
				http.Error(w, "Ошибка чтения данных", http.StatusInternalServerError)
				return
			}
			items = append(items, item)
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(items)
		
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	
}
