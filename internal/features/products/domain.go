package products

import "time"

// Product is the domain model for the products feature.
// It is decoupled from the database implementation and uses clean Go types.
type Product struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Price     string    `json:"price"`
	RecStatus string    `json:"rec_status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
