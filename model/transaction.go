package model

import "time"

type Transaction struct {
	ID          int                 `json:"id" db:"id"`
	TotalAmount int                 `json:"total_amount" db:"total_amount"`
	CreatedAt   time.Time           `json:"created_at" db:"created_at"`
	Details     []TransactionDetail `json:"details" db:"-"`
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
type ProductTerlaris struct {
	Nama       string `json:"nama" db:"product_name"`
	QtyTerjual int    `json:"qty_terjual" db:"total_qty"`
}

type Report struct {
	TotalRevenue    int             `json:"total_revenue" db:"total_revenue"`
	TotalTransaks   int             `json:"total_transaksi" db:"total_transaksi"`
	ProductTerlaris ProductTerlaris `json:"product_terlaris"`
}
