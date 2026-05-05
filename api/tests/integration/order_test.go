package integration_test

import (
	"net/http"
	"testing"
	"yarik-varit/api/models"
)

func TestGetOrders_WithoutAuth(t *testing.T) {
	res := get("/orders", "")

	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("ожидали 401, получили %d", res.StatusCode)
	}
}

func TestGetOrders_AsAdmin(t *testing.T) {
	res := get("/orders", adminToken)

	if res.StatusCode != http.StatusOK {
		t.Fatalf("ожидали 200, получили %d", res.StatusCode)
	}

	var orders []models.Order
	decode(res, &orders)
	// orders может быть пустым — это нормально
}

func TestGetOrders_FilterByStatus(t *testing.T) {
	res := get("/orders?status=pending", adminToken)

	if res.StatusCode != http.StatusOK {
		t.Fatalf("ожидали 200, получили %d", res.StatusCode)
	}

	var orders []models.Order
	decode(res, &orders)
	for _, o := range orders {
		if o.Status != models.StatusPending {
			t.Errorf("ожидали статус pending, получили %s", o.Status)
		}
	}
}

func TestCreateOrder_EmptyItems(t *testing.T) {
	res := post("/orders", adminToken, map[string]interface{}{
		"items": []interface{}{},
	})

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("ожидали 400, получили %d", res.StatusCode)
	}
}

func TestCreateOrder_InvalidMenuID(t *testing.T) {
	res := post("/orders", adminToken, map[string]interface{}{
		"items": []map[string]interface{}{
			{"menu_id": 999999, "quantity": 1},
		},
	})

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("ожидали 400, получили %d", res.StatusCode)
	}
}

func TestCreateOrder_Success(t *testing.T) {
	// Берём реальный menu_id из БД
	var menuID int
	err := testDB.QueryRow("SELECT id FROM menu LIMIT 1").Scan(&menuID)
	if err != nil {
		t.Skip("В меню нет позиций — пропускаем тест")
	}

	res := post("/orders", adminToken, map[string]interface{}{
		"items": []map[string]interface{}{
			{"menu_id": menuID, "quantity": 2},
		},
	})

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("ожидали 201, получили %d: %s", res.StatusCode, bodyString(res))
	}

	var order models.Order
	decode(res, &order)

	if order.Total <= 0 {
		t.Error("total должен быть больше 0")
	}
	if len(order.Item) == 0 {
		t.Error("в заказе должны быть позиции")
	}

	// Чистим тестовый заказ
	testDB.Exec("DELETE FROM order_items WHERE order_id = $1", order.ID)
	testDB.Exec("DELETE FROM orders WHERE id = $1", order.ID)
}

func TestGetOrderByID_NotFound(t *testing.T) {
	res := get("/orders/999999", adminToken)

	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("ожидали 404, получили %d", res.StatusCode)
	}
}

func TestUpdateOrderStatus_InvalidStatus(t *testing.T) {
	res, _ := doRequest("PATCH", "/orders/1/status", adminToken, map[string]string{
		"status": "flying",
	})

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("ожидали 400, получили %d", res.StatusCode)
	}
}
