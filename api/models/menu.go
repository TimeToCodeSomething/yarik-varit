package models

type MenuItem struct {
	Name  string `json:"name"`
	Price int `json:"price"`
	Vol   int `json:"vol"`
}