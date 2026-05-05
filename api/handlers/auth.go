package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"
	"yarik-varit/api/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func jwtSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("Переменная JWT_SECRET не задана")
	}
	return []byte(secret)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ошибка чтения запроса", http.StatusBadRequest)
		return
	}

	if len(req.Username) < 3 {
		http.Error(w, "Логин должен быть не короче 3 символов", http.StatusBadRequest)
		return
	}
	if len(req.Password) < 6 {
		http.Error(w, "Пароль должен быть не короче 6 символов", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Ошибка хэширования пароля", http.StatusInternalServerError)
		return
	}

	var user models.User
	err = h.db.QueryRow(
		"INSERT INTO users (username, password_hash, role) VALUES ($1, $2, $3) RETURNING id, username, role",
		req.Username, string(hash), models.RoleClient,
	).Scan(&user.ID, &user.Username, &user.Role)

	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			http.Error(w, "Логин уже занят", http.StatusConflict)
			return
		}
		http.Error(w, "Ошибка создания пользователя", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ошибка чтения запроса", http.StatusBadRequest)
		return
	}

	var user models.User
	err := h.db.QueryRow(
		"SELECT id, username, password_hash, role FROM users WHERE username = $1", req.Username,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role)

	if err == sql.ErrNoRows {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Ошибка БД", http.StatusInternalServerError)
		return
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
		return
	}

	// Генерируем JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Токен живет 3 дня
	})

	tokenString, err := token.SignedString(jwtSecret())
	if err != nil {
		http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
