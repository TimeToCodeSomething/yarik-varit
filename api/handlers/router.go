package handlers

import (
	"database/sql"
	"net/http"
	"yarik-varit/api/models"
)

func NewRouter(db *sql.DB) *http.ServeMux {
	h := NewHandler(db)
	mux := http.NewServeMux()

	// Публичные эндпоинты
	mux.HandleFunc("POST /login", h.Login)
	mux.HandleFunc("POST /register", h.Register)
	mux.HandleFunc("GET /menu", h.GetMenuItems)
	mux.HandleFunc("GET /menu/categories", h.GetMenuByCategories)

	// Роли
	adminOnly := []models.Role{models.RoleAdmin}
	staffRoles := []models.Role{models.RoleAdmin, models.RoleBarista}
	anyUser := []models.Role{models.RoleAdmin, models.RoleBarista, models.RoleClient}

	// Menu endpoints
	mux.HandleFunc("POST /menu", AuthMiddleware(adminOnly, h.CreateMenuItem))
	mux.HandleFunc("PUT /menu/{id}", AuthMiddleware(adminOnly, h.UpdateMenuItem))
	mux.HandleFunc("DELETE /menu/{id}", AuthMiddleware(adminOnly, h.DeleteMenuItem))

	// Orders endpoints
	mux.HandleFunc("GET /orders", AuthMiddleware(staffRoles, h.GetOrders))
	mux.HandleFunc("GET /orders/{id}", AuthMiddleware(anyUser, h.GetOrderByID))
	mux.HandleFunc("POST /orders", AuthMiddleware(anyUser, h.CreateOrder))
	mux.HandleFunc("PATCH /orders/{id}/status", AuthMiddleware(staffRoles, h.UpdateOrderStatus))
	mux.HandleFunc("DELETE /orders/{id}", AuthMiddleware(staffRoles, h.DeleteOrder))

	return mux
}
