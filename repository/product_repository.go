package repository

import (
	"category-crud/model"
	"category-crud/model/dto"
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v9"
)

type ProductRepository struct {
	db      *sql.DB
	builder *goqu.Database
}

func NewProductRepository(db *sql.DB, builder *goqu.Database) *ProductRepository {
	return &ProductRepository{
		db:      db,
		builder: builder,
	}
}

func (repo *ProductRepository) GetAll(filter *dto.ProductFilterRequest) ([]model.Product, error) {
	// Get all products
	var products []model.Product
	queryRaw := repo.builder.
		From("products").
		Select("id", "name", "price", "stock")
	if filter.Name != "" {
		queryRaw = queryRaw.Where(goqu.I("name").ILike("%" + filter.Name + "%"))
	}
	err := queryRaw.
		ScanStructs(&products)
	if err != nil {
		return nil, err
	}

	productIDs := make([]int, len(products))

	for _, product := range products {
		productIDs = append(productIDs, product.ID)
	}

	if len(productIDs) == 0 {
		return products, nil
	}

	// Get all categories for these products
	categoryQuery, _, err := repo.builder.From(goqu.T("product_categories").As("pc")).
		Select(
			goqu.I("pc.product_id"),
			goqu.I("c.id"),
			goqu.I("c.name"),
		).
		Join(
			goqu.T("categories").As("c"),
			goqu.On(goqu.Ex{"pc.category_id": goqu.I("c.id")}),
		).
		Where(goqu.I("pc.product_id").In(productIDs)). // productIds is your slice
		Order(
			goqu.I("pc.product_id").Asc(),
			goqu.I("c.name").Asc(),
		).ToSQL()

	catRows, err := repo.db.Query(categoryQuery)
	if err != nil {
		return nil, err
	}

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

func GenerateInsertProductCategoriesQuery(builder *goqu.Database, product *dto.ProductRequest) string {
	records := make([]goqu.Record, 0, len(product.Categories))

	for _, categoryID := range product.Categories {
		records = append(records, goqu.Record{
			"product_id":  product.ID,
			"category_id": categoryID,
		})
	}

	queryCategoryInsert, _, err := builder.Insert("product_categories").Rows(records).ToSQL()

	if err != nil {
		return ""
	}

	return queryCategoryInsert

}

func (repo *ProductRepository) Create(product *dto.ProductRequest) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query, _, err := repo.builder.Insert("products").Rows(
		goqu.Record{
			"name":  product.Name,
			"price": product.Price,
			"stock": product.Stock,
		},
	).Returning("id").ToSQL()
	err = tx.QueryRow(query).Scan(&product.ID)
	// Batch insert categories
	if len(product.Categories) > 0 {
		records := make([]goqu.Record, 0, len(product.Categories))

		for _, categoryID := range product.Categories {
			records = append(records, goqu.Record{
				"product_id":  product.ID,
				"category_id": categoryID,
			})
		}

		queryCategoryInsert := GenerateInsertProductCategoriesQuery(repo.builder, product)
		_, err = tx.Exec(queryCategoryInsert)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetByID - ambil produk by ID
func (repo *ProductRepository) GetByID(id int) (*model.Product, error) {
	var product model.Product
	result, err := repo.builder.
		From("products").
		Select("id", "name", "price", "stock").
		Where(goqu.Ex{"id": id}).
		ScanStruct(&product)

	if !result {
		return nil, errors.New("produk tidak ditemukan")
	}

	if err != nil {
		return nil, err
	}
	var categories []model.Category
	err = repo.builder.From(goqu.T("categories").As("c")).
		Select(
			goqu.I("c.id"),
			goqu.I("c.name"),
		).
		Join(
			goqu.T("product_categories").As("pc"),
			goqu.On(goqu.Ex{"c.id": goqu.I("pc.category_id")}),
		).
		Where(goqu.Ex{"pc.product_id": id}).
		Order(goqu.I("c.name").Asc()).
		ScanStructs(&categories)

	if err != nil {
		return nil, err
	}

	product.Categories = categories

	return &product, nil
}

func (repo *ProductRepository) Update(product *dto.ProductRequest) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query, _, err := repo.builder.Update("products").Set(
		goqu.Record{
			"name":  product.Name,
			"price": product.Price,
			"stock": product.Stock,
		}).
		Where(goqu.Ex{"id": product.ID}).ToSQL()

	result, err := tx.Exec(query)
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
	deleteQuery, _, err := repo.builder.
		Delete("product_categories").
		Where(goqu.Ex{"product_id": product.ID}).ToSQL()

	_, err = tx.Exec(deleteQuery)
	if err != nil {
		return err
	}

	// Batch insert new categories
	if len(product.Categories) > 0 {
		insertQuery := GenerateInsertProductCategoriesQuery(repo.builder, product)

		_, err = tx.Exec(insertQuery)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (repo *ProductRepository) Delete(id int) error {
	query, _, err := repo.builder.Delete("products").Where(goqu.Ex{"id": id}).ToSQL()
	result, err := repo.db.Exec(query)
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
