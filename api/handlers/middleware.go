package handlers

import (
	"context"
	"net/http"
	"strings"
	"yarik-varit/api/models"

	"github.com/golang-jwt/jwt/v5"
)

// Ключ для хранения данных пользователя в контексте запроса
type contextKey string

const userContextKey contextKey = "userCtx"

type UserContextData struct {
	UserID int
	Role   string
}

// AuthMiddleware принимает список разрешенных ролей и следующий хендлер
func AuthMiddleware(allowedRoles []models.Role, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
			return
		}

		// Ожидаем формат "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Неверный формат заголовка Authorization", http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]

		// Парсим и валидируем токен
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Недействительный токен", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Ошибка чтения токена", http.StatusInternalServerError)
			return
		}

		userRole := claims["role"].(string)
		userID := int(claims["user_id"].(float64))

		// Проверяем, есть ли роль пользователя в списке разрешенных
		roleAllowed := false
		for _, allowed := range allowedRoles {
			if string(allowed) == userRole {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			http.Error(w, "Недостаточно прав", http.StatusForbidden)
			return
		}

		// Кладем данные пользователя в контекст, чтобы их можно было достать в самом хендлере (например, чтобы понять, чей это заказ)
		ctx := context.WithValue(r.Context(), userContextKey, UserContextData{
			UserID: userID,
			Role:   userRole,
		})

		// Передаем управление дальше по цепочке
		next(w, r.WithContext(ctx))
	}
}
