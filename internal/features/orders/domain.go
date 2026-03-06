package orders

import "time"

// Order represents the order header
type Order struct {
	ID           string      `json:"id"`
	CustomerName string      `json:"customer_name"`
	TotalAmount  float64     `json:"total_amount"`
	CreatedAt    time.Time   `json:"created_at"`
	Items        []OrderItem `json:"items,omitempty"`
}

// OrderItem represents each product in an order, enriched with product details
type OrderItem struct {
	ID          int32   `json:"id"`
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"` // From JOIN
	Quantity    int32   `json:"quantity"`
	Price       float64 `json:"price"` // Historical price at the time of order
}
