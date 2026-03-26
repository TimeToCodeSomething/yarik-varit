package models

type MenuItem struct {
	ID 			int			`json:"id"`
	Name 		string		`json:"name"`
	Price 		float64			`json:"price"`
	Vol 		*int		`json:"vol"`
	Category	string		`json:"category"`
}