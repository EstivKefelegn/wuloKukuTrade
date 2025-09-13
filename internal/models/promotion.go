package models

import (
	"time"

	"github.com/govalues/decimal"
)

type Promotion struct {
	ID          string          `josn:"id,omiempty" db:"id,omitempty"`
	Description string          `json:"description,omitempty" db:"description,omiempty"`
	Discount    decimal.Decimal `json:"discount,omitempty" db:"discount,omitempty"`
}

// for product and Promotion many to many relationship
type PromotionProduct struct {
	ID          string     `json:"id" db:"id"`
	PromotionID string     `json:"promotion_id" db:"promotion_id"`
	ProductID   string     `json:"product_id" db:"product_id"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	Promotion   *Promotion `json:"promotion,omitempty" db:"-"`
	Product     *Product   `json:"product,omitempty" db:"-"`
}
