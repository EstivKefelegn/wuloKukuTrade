package models

import "time"

type Cart struct {
	ID        string    `json:"id,omitempty" db:"id,omitempty"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CartItem struct {
	ID        int      `json:"id,omitempty" db:"id,omitempty"`
	CartID    string   `json:"cart_id" db:"cart_id"`
	ProductID string   `json:"product_id" db:"product_id"`
	Quantity  int      `json:"quantity" db:"quantity"`
	Cart      *Cart    `json:"-"`
	Product   *Product `json:"-"`
}
