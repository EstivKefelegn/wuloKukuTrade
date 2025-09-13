package models

import (
	"time"
)

type Product struct {
	ID           int         `json:"id,omitempty" db:"id,omitempty"`
	Breed        string      `json:"breed" db:"breed"`
	AgeWeek      int         `json:"age_week" db:"age_week"`
	PricePerHen  float32     `json:"price_per_hen" db:"price_per_hen"`
	IsVaccinated bool        `json:"is_vaccinated" db:"is_vaccinated"`
	LastUpdate   time.Time   `json:"last_update" db:"last_update"`
	CollectionID int         `json:"collection_id" db:"collection_id"`
	Collection   *Collection `json:"-"`
}

type Review struct {
	ID          int       `json:"id,omitempty" db:"id,omitempty"`
	Product_id  int       `json:"product_id" db:"product_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Date        time.Time `json:"date" db:"date"`
	Product     *Product  `json:"-"`
}
