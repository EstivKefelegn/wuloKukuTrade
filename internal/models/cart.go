package models

import "time"

type Cart struct {
	ID        int       `json:"id,omitempty" db:"id,omitempty"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CartItem struct {
	ID        int      `json:"id,omitempty" db:"id,omitempty"`
	CartID    int      `json:"cart_id" db:"cart_id"`
	ProductID int      `json:"product_id" db:"product_id"`
	Quantity  int      `json:"quantity" db:"quantity"`
	Cart      *Cart    `json:"-"`
	Product   *Product `json:"-"`
}
