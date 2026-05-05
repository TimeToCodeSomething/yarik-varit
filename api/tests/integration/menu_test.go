package integration_test

import (
	"fmt"
	"net/http"
	"testing"
	"yarik-varit/api/models"
)

func TestGetMenu_Public(t *testing.T) {
	res := get("/menu", "")

	if res.StatusCode != http.StatusOK {
		t.Fatalf("ожидали 200, получили %d", res.StatusCode)
	}

	var items []models.MenuItem
	decode(res, &items)
	if len(items) == 0 {
		t.Error("меню пустое — добавьте позиции в БД")
	}
}

func TestGetMenu_FilterByCategory(t *testing.T) {
	res := get("/menu?category=espresso", "")

	if res.StatusCode != http.StatusOK {
		t.Fatalf("ожидали 200, получили %d", res.StatusCode)
	}

	var items []models.MenuItem
	decode(res, &items)
	for _, item := range items {
		if item.Category != "espresso" {
			t.Errorf("ожидали категорию espresso, получили %s", item.Category)
		}
	}
}

func TestCreateMenuItem_WithoutAuth(t *testing.T) {
	res := post("/menu", "", map[string]interface{}{
		"name": "Тест", "price": 100, "vol": 200, "category": "espresso",
	})

	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("ожидали 401, получили %d", res.StatusCode)
	}
}

func TestCreateMenuItem_AsAdmin(t *testing.T) {
	res := post("/menu", adminToken, map[string]interface{}{
		"name": "Тестовый напиток", "price": 250, "vol": 150, "category": "espresso",
	})

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("ожидали 201, получили %d: %s", res.StatusCode, bodyString(res))
	}

	var item models.MenuItem
	decode(res, &item)

	// Чистим тестовую запись
	testDB.Exec("DELETE FROM menu WHERE id = $1", item.ID)
}

func TestCreateMenuItem_AsClient(t *testing.T) {
	// Регистрируем временного клиента
	username := fmt.Sprintf("client_%d", uniqueCounter())
	post("/register", "", map[string]string{"username": username, "password": "pass123456"})
	defer testDB.Exec("DELETE FROM users WHERE username = $1", username)

	clientToken := mustLogin(username, "pass123456")

	res := post("/menu", clientToken, map[string]interface{}{
		"name": "Тест", "price": 100, "vol": 200, "category": "espresso",
	})

	if res.StatusCode != http.StatusForbidden {
		t.Fatalf("ожидали 403, получили %d", res.StatusCode)
	}
}
