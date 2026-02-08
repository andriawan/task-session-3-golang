package model

import "time"

type Transaction struct {
	ID          int                 `json:"id"`
	TotalAmount int                 `json:"total_amount"`
	CreatedAt   time.Time           `json:"created_at"`
	Details     []TransactionDetail `json:"details"`
}

type TransactionDetail struct {
	ID            int    `json:"id" db:"id"`
	TransactionID int    `json:"transaction_id" db:"transaction_id"`
	ProductID     int    `json:"product_id" db:"product_id"`
	ProductName   string `json:"product_name,omitempty" db:"product_name"`
	Quantity      int    `json:"quantity" db:"quantity"`
	Subtotal      int    `json:"subtotal" db:"subtotal"`
}

type CheckoutItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}
