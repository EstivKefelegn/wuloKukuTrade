package models

type Collection struct {
	ID                int    `json:"id,omitempty" db:"id,omitempty"`
	Title             string `json:"title" db:"title"`
	FeaturedProductID *int    `json:"featured_product_id" db:"featured_product_id"`
}
