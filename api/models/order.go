package models

import (
	"time"
)

type OrderStatus string

const (
	StatusPending  OrderStatus = "pending"
	StatusReady    OrderStatus = "ready"
	StatusDone     OrderStatus = "done"
	StatusCanceled OrderStatus = "canceled"
)

type Order struct {
	ID     int         `json:"id"`
	Time   time.Time   `json:"time"`
	Status OrderStatus `json:"status"`
	UserID *int        `json:"user_id,omitempty"`
	Total  float64     `json:"total"`
	Item   []OrderItem `json:"items,omitempty"`
}

type OrderItem struct {
	ID       int     `json:"id"`
	OrderID  int     `json:"order_id"`
	MenuID   int     `json:"menu_id"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	Name     string  `json:"name"`
}
