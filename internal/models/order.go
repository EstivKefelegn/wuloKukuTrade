package models

import (
	"time"

	"github.com/govalues/decimal"
)

type PaymentStatus string

const (
	Pending PaymentStatus = "pending"
	Failed  PaymentStatus = "failed"
	Success PaymentStatus = "success"
)

type Orders struct {
	ID            int           `json:"id,omitempty" db:"id,omitempty"`
	PlacedAt      time.Time     `json:"placed_at" db:"placed_at"`
	PaymentStatus PaymentStatus `json:"payment_status" db:"payment_status"`
	UserId        int           `json:"user_id" db:"user_id"`
	User          *Users        `json:"-"`
}

type OrderItem struct {
	ID        int             `json:"id,omitempty" db:"id,omitempty"`
	OrderID   int             `json:"order_id" db:"order_id"`
	ProductID int             `json:"product_id" db:"product_id"`
	Quantity  int             `json:"quantity" db:"quantity"`
	Price     decimal.Decimal `json:"price" db:"price"`
	Order     *Orders         `json:"-"`
	Product   *Product        `json:"-"`
}

func (p PaymentStatus) isValid() bool {
	switch p {
	case Pending, Failed, Success:
		return true
	}

	return false
}
