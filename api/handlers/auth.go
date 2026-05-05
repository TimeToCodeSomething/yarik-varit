package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
	"yarik-varit/api/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// В идеале секрет должен браться из .env (os.Getenv("JWT_SECRET"))
var jwtSecret = []byte("super-secret-yarik-key")

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
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

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
