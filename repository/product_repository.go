package repository

import (
	"category-crud/model"
	"category-crud/model/dto"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (repo *ProductRepository) GetAll() ([]model.Product, error) {
	// Get all products
	query := "SELECT id, name, price, stock FROM products ORDER BY id"
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]model.Product, 0)
	productIDs := make([]int, 0)

	for rows.Next() {
		var p model.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
		if err != nil {
			return nil, err
		}
		p.Categories = []model.Category{} // Initialize empty slice
		products = append(products, p)
		productIDs = append(productIDs, p.ID)
	}

	if len(productIDs) == 0 {
		return products, nil
	}

	// Get all categories for these products
	categoryQuery := `
		SELECT pc.product_id, c.id, c.name
		FROM product_categories pc
		JOIN categories c ON pc.category_id = c.id
		WHERE pc.product_id = ANY($1)
		ORDER BY pc.product_id, c.name
	`

	catRows, err := repo.db.Query(categoryQuery, pq.Array(productIDs))
	if err != nil {
		return nil, err
	}
	defer catRows.Close()

	// Map categories to products
	productMap := make(map[int]*model.Product)
	for i := range products {
		productMap[products[i].ID] = &products[i]
	}

	for catRows.Next() {
		var productID, categoryID int
		var categoryName string

		err := catRows.Scan(&productID, &categoryID, &categoryName)
		if err != nil {
			return nil, err
		}

		if product, exists := productMap[productID]; exists {
			product.Categories = append(product.Categories, model.Category{
				ID:   categoryID,
				Name: categoryName,
			})
		}
	}

	return products, nil
}

func (repo *ProductRepository) Create(product *dto.ProductRequest) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := "INSERT INTO products (name, price, stock) VALUES ($1, $2, $3) RETURNING id"
	err = tx.QueryRow(query, product.Name, product.Price, product.Stock).Scan(&product.ID)
	// Batch insert categories
	if len(product.Categories) > 0 {
		// Build bulk insert query
		valueStrings := []string{}
		valueArgs := []any{}

		for i, categoryID := range product.Categories {
			valueStrings = append(valueStrings, fmt.Sprintf("($1, $%d)", i+2))
			valueArgs = append(valueArgs, categoryID)
		}

		queryCategoryInsert := fmt.Sprintf(
			"INSERT INTO product_categories (product_id, category_id) VALUES %s",
			strings.Join(valueStrings, ","),
		)

		// Prepend product ID to args
		valueArgs = append([]any{product.ID}, valueArgs...)

		_, err = tx.Exec(queryCategoryInsert, valueArgs...)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetByID - ambil produk by ID
func (repo *ProductRepository) GetByID(id int) (*model.Product, error) {
	query := "SELECT id, name, price, stock FROM products WHERE id = $1"

	var p model.Product
	err := repo.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
	if err == sql.ErrNoRows {
		return nil, errors.New("produk tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}

	// Get categories for this product
	categoryQuery := `
		SELECT c.id, c.name
		FROM categories c
		JOIN product_categories pc ON c.id = pc.category_id
		WHERE pc.product_id = $1
		ORDER BY c.name
	`

	rows, err := repo.db.Query(categoryQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	p.Categories = []model.Category{}
	for rows.Next() {
		var cat model.Category
		err := rows.Scan(&cat.ID, &cat.Name)
		if err != nil {
			return nil, err
		}
		p.Categories = append(p.Categories, cat)
	}

	return &p, nil
}

func (repo *ProductRepository) Update(product *dto.ProductRequest) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := "UPDATE products SET name = $1, price = $2, stock = $3 WHERE id = $4"
	result, err := tx.Exec(query, product.Name, product.Price, product.Stock, product.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}

	// Delete existing category relationships
	deleteQuery := "DELETE FROM product_categories WHERE product_id = $1"
	_, err = tx.Exec(deleteQuery, product.ID)
	if err != nil {
		return err
	}

	// Batch insert new categories
	if len(product.Categories) > 0 {
		valueStrings := make([]string, 0, len(product.Categories))
		valueArgs := make([]any, 0, len(product.Categories)+1)
		valueArgs = append(valueArgs, product.ID)

		for i, categoryID := range product.Categories {
			valueStrings = append(valueStrings, fmt.Sprintf("($1, $%d)", i+2))
			valueArgs = append(valueArgs, categoryID)
		}

		insertQuery := fmt.Sprintf(
			"INSERT INTO product_categories (product_id, category_id) VALUES %s",
			strings.Join(valueStrings, ","),
		)

		_, err = tx.Exec(insertQuery, valueArgs...)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (repo *ProductRepository) Delete(id int) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}

	return err
}
