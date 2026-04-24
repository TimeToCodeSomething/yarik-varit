package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"yarik-varit/api/models"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type OrderRequest struct {
	Items []struct {
		MenuID   int `json:"menu_id"`
		Quantity int `json:"quantity"`
	} `json:"items"`
}

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query("SELECT id, tm, status, total FROM orders ORDER BY tm DESC")
	if err != nil {
		http.Error(w, "Ошибка выполнения запроса", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.Time, &o.Status, &o.Total); err != nil {
			http.Error(w, "Ошибка чтения данных", http.StatusInternalServerError)
			return
		}
		orders = append(orders, o)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Ошибка итерации", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(orders)
}

func (h *Handler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	idInt, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Неверный id", http.StatusBadRequest)
		return
	}

	var order models.Order
	err = h.db.QueryRow(
		"SELECT id, tm, status, total FROM orders WHERE id = $1", idInt,
	).Scan(&order.ID, &order.Time, &order.Status, &order.Total)
	if err == sql.ErrNoRows {
		http.Error(w, "Заказ не найден", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Ошибка запроса", http.StatusInternalServerError)
		return
	}

	rows, err := h.db.Query(
		"SELECT id, order_id, menu_id, quantity, price, name FROM order_items WHERE order_id = $1",
		idInt,
	)
	if err != nil {
		http.Error(w, "Ошибка получения позиций", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.MenuID, &item.Quantity, &item.Price, &item.Name); err != nil {
			http.Error(w, "Ошибка чтения позиций", http.StatusInternalServerError)
			return
		}
		order.Item = append(order.Item, item)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Ошибка итерации позиций", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(order)
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ошибка чтения запроса", http.StatusBadRequest)
		return
	}

	if len(req.Items) == 0 {
		http.Error(w, "Заказ не может быть пустым", http.StatusBadRequest)
		return
	}
	for _, item := range req.Items {
		if item.MenuID <= 0 {
			http.Error(w, "Неверный menu_id", http.StatusBadRequest)
			return
		}
		if item.Quantity <= 0 {
			http.Error(w, "Количество должно быть больше 0", http.StatusBadRequest)
			return
		}
	}

	tx, err := h.db.Begin()
	if err != nil {
		http.Error(w, "Ошибка транзакции", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Сначала считаем total и собираем данные
	var total float64
	itemsData := make([]struct {
		MenuID   int
		Quantity int
		Price    float64
		Name     string
	}, len(req.Items))

	for i, item := range req.Items {
		var price float64
		var name string
		err := tx.QueryRow(
			"SELECT price, name FROM menu WHERE id = $1", item.MenuID,
		).Scan(&price, &name)
		if err == sql.ErrNoRows {
			http.Error(w, "Позиция меню не найдена", http.StatusBadRequest)
			return
		}
		if err != nil {
			http.Error(w, "Ошибка получения цены", http.StatusInternalServerError)
			return
		}

		itemsData[i].MenuID = item.MenuID
		itemsData[i].Quantity = item.Quantity
		itemsData[i].Price = price
		itemsData[i].Name = name
		total += price * float64(item.Quantity)
	}

	// Создаём заказ с total
	var order models.Order
	err = tx.QueryRow(
		"INSERT INTO orders (status, total) VALUES ($1, $2) RETURNING id, tm",
		models.StatusPending, total,
	).Scan(&order.ID, &order.Time)
	if err != nil {
		http.Error(w, "Ошибка создания заказа", http.StatusInternalServerError)
		return
	}
	order.Status = models.StatusPending
	order.Total = total

	// Добавляем items
	for _, itemData := range itemsData {
		var oi models.OrderItem
		err = tx.QueryRow(
			"INSERT INTO order_items (order_id, menu_id, quantity, price, name) VALUES ($1, $2, $3, $4, $5) RETURNING id, order_id, menu_id, quantity, price, name",
			order.ID, itemData.MenuID, itemData.Quantity, itemData.Price, itemData.Name,
		).Scan(&oi.ID, &oi.OrderID, &oi.MenuID, &oi.Quantity, &oi.Price, &oi.Name)
		if err != nil {
			http.Error(w, "Ошибка добавления позиции", http.StatusInternalServerError)
			return
		}
		order.Item = append(order.Item, oi)
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Ошибка подтверждения заказа", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func (h *Handler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	idInt, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Неверный id", http.StatusBadRequest)
		return
	}

	var req struct {
		Status models.OrderStatus `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ошибка чтения запроса", http.StatusBadRequest)
		return
	}

	valid := false
	for _, s := range []models.OrderStatus{models.StatusPending, models.StatusReady, models.StatusDone, models.StatusCanceled} {
		if req.Status == s {
			valid = true
			break
		}
	}
	if !valid {
		http.Error(w, "Неверный статус", http.StatusBadRequest)
		return
	}

	result, err := h.db.Exec("UPDATE orders SET status = $1 WHERE id = $2", req.Status, idInt)
	if err != nil {
		http.Error(w, "Ошибка обновления статуса", http.StatusInternalServerError)
		return
	}

	// FIX: обрабатываем ошибку RowsAffected
	affected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Ошибка проверки обновления", http.StatusInternalServerError)
		return
	}
	if affected == 0 {
		http.Error(w, "Заказ не найден", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
