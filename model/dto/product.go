package dto

type ProductRequest struct {
	ID         int    `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	Price      int    `json:"price" db:"price"`
	Stock      int    `json:"stock" db:"stock"`
	Categories []int  `json:"categories" db:"-"`
}

type ProductFilterRequest struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	IDs  []int  `json:"ids"`
}
