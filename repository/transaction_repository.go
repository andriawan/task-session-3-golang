package repository

import (
	"category-crud/model"
	"category-crud/model/dto"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
)

type TransactionRepository struct {
	db          *sql.DB
	productRepo *ProductRepository
	builder     *goqu.Database
}

type TransactionResult struct {
	ID        uint64    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
}

func NewTransactionRepository(db *sql.DB, builder *goqu.Database, productRepo *ProductRepository) *TransactionRepository {
	return &TransactionRepository{
		db:          db,
		productRepo: productRepo,
		builder:     builder,
	}
}

func (repo *TransactionRepository) CreateTransaction(items []model.CheckoutItem) (*model.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]model.TransactionDetail, 0, len(items))
	productID := make([]int, 0, len(items))

	// get product id mapping
	for _, item := range items {
		productID = append(productID, item.ProductID)
	}

	products, err := repo.productRepo.GetAll(&dto.ProductFilterRequest{
		IDs: productID,
	})

	productMap := make(map[int]*model.Product)

	for _, product := range products {
		productMap[product.ID] = &product
	}

	// map total amount and update stock
	for _, item := range items {
		product, ok := productMap[item.ProductID]
		if !ok {
			return nil, err
		}

		totalAmount += product.Price * item.Quantity
		details = append(details, model.TransactionDetail{
			ProductID:   product.ID,
			ProductName: product.Name,
			Quantity:    item.Quantity,
			Subtotal:    product.Price * item.Quantity,
		})
		// update stock
		product.Stock -= item.Quantity
	}

	// insert total Amount
	var result TransactionResult

	_, err = repo.builder.Insert("transactions").Rows(
		goqu.Record{
			"total_amount": totalAmount,
		},
	).Returning("id", "created_at").Executor().ScanStruct(&result)

	if err != nil {
		return nil, err
	}

	// insert details
	var insertedDetails []model.TransactionDetail
	detailRecords := make([]goqu.Record, 0, len(details))
	for _, detail := range details {
		detailRecords = append(detailRecords, goqu.Record{
			"transaction_id": result.ID,
			"product_id":     detail.ProductID,
			"quantity":       detail.Quantity,
			"subtotal":       detail.Subtotal,
		})
	}

	err = repo.builder.Insert("transaction_details").Rows(
		detailRecords,
	).
		Returning(goqu.Star()).
		Executor().ScanStructs(&insertedDetails)

	if err != nil {
		return nil, err
	}

	// update product stock
	for _, product := range products {
		_, err = repo.builder.Update("products").Set(goqu.Record{
			"stock": productMap[product.ID].Stock,
		}).Where(goqu.Ex{
			"id": product.ID,
		}).Executor().Exec()
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &model.Transaction{
		ID:          int(result.ID),
		CreatedAt:   result.CreatedAt,
		TotalAmount: totalAmount,
		Details:     insertedDetails,
	}, nil
}
