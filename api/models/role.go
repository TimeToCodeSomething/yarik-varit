package models

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleBarista Role = "barista"
	RoleClient  Role = "client"
)

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"` // Пароль не должен отдаваться в JSON
	Role         Role   `json:"role"`
}
