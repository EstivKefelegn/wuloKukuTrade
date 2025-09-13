package models

type Collection struct {
	ID                string   `json:"id,omitempty" db:"id,omitempty"`
	Title             string   `json:"title" db:"title"`
	FeaturedProductID string   `json:"featured_product_id" db:"featured_product_id"`
	FeaturedProduct   *Product `json:"-"`
}
