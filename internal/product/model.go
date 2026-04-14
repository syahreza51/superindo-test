package product

import "time"

type Type string

const (
	TypeSayuran Type = "Sayuran"
	TypeProtein Type = "Protein"
	TypeBuah    Type = "Buah"
	TypeSnack   Type = "Snack"
)

func (t Type) Valid() bool {
	switch t {
	case TypeSayuran, TypeProtein, TypeBuah, TypeSnack:
		return true
	default:
		return false
	}
}

type Product struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      Type      `json:"type"`
	Price     int64     `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

