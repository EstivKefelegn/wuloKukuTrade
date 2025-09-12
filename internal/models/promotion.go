package models

import "github.com/govalues/decimal"

type Promotion struct {
	ID          int             `josn:"id,omiempty" db:"id,omitempty"`
	Description string          `json:"description,omitempty" db:"description,omiempty"`
	Discount    decimal.Decimal `json:"discount,omitempty" db:"discount,omitempty"`
}
