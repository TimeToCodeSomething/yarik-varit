package models

import (
	"time"
)

type OrderStatus string

const (
	StatusPending OrderStatus = "pending"
	StatusReady   OrderStatus = "ready"
	StatusDone    OrderStatus = "done"
)

type Order struct {
	ID     int         `json:"id"`
	Time   time.Time   `json:"time"`
	Status OrderStatus `json:"status"`
}

type OrderItem struct {
	ID       int `json:"id"`
	OrderID  int `json:"order_id"`
	MenuID   int `json:"menu_id"`
	Quantity int `json:"quantity"`
}
