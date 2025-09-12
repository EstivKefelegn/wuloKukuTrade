package models

import (
	"time"
)

type Product struct {
	ID           int       `json:"id,omitempty" db:"id,omiempty"`
	Title        string    `json:"title" db:"title"`
	Inventory    int       `json:"inventory" db:"inventory"`
	LastUpdate   time.Time `json:"last_update" db:"last_update"`
	CollectionID int       `json:"collection_id" db:"collection_id"`
	Collection   *Collection
}


