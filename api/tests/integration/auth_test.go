package integration_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

func TestLogin_Success(t *testing.T) {
	res := post("/login", "", map[string]string{
		"username": os.Getenv("TEST_ADMIN_USER"),
		"password": os.Getenv("TEST_ADMIN_PASS"),
	})

	if res.StatusCode != http.StatusOK {
		t.Fatalf("ожидали 200, получили %d: %s", res.StatusCode, bodyString(res))
	}

	var data map[string]string
	decode(res, &data)
	if data["token"] == "" {
		t.Fatal("токен пустой в ответе")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	res := post("/login", "", map[string]string{
		"username": os.Getenv("TEST_ADMIN_USER"),
		"password": "wrongpassword",
	})

	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("ожидали 401, получили %d", res.StatusCode)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	res := post("/login", "", map[string]string{
		"username": "nosuchuser",
		"password": "password",
	})

	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("ожидали 401, получили %d", res.StatusCode)
	}
}

func TestRegister_Success(t *testing.T) {
	username := fmt.Sprintf("testuser_%d", uniqueCounter())
	res := post("/register", "", map[string]string{
		"username": username,
		"password": "testpass123",
	})

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("ожидали 201, получили %d: %s", res.StatusCode, bodyString(res))
	}

	// Чистим тестового пользователя
	testDB.Exec("DELETE FROM users WHERE username = $1", username)
}

func TestRegister_DuplicateUsername(t *testing.T) {
	res := post("/register", "", map[string]string{
		"username": os.Getenv("TEST_ADMIN_USER"),
		"password": "somepassword",
	})

	if res.StatusCode != http.StatusConflict {
		t.Fatalf("ожидали 409, получили %d", res.StatusCode)
	}
}

func TestRegister_ShortPassword(t *testing.T) {
	res := post("/register", "", map[string]string{
		"username": "newuser",
		"password": "123",
	})

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("ожидали 400, получили %d", res.StatusCode)
	}
}

func TestRegister_ShortUsername(t *testing.T) {
	res := post("/register", "", map[string]string{
		"username": "ab",
		"password": "password123",
	})

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("ожидали 400, получили %d", res.StatusCode)
	}
}

var counter int

func uniqueCounter() int {
	counter++
	return counter
}
