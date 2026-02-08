package repository

import (
	"category-crud/model"
	"category-crud/model/dto"
	"database/sql"
	"errors"
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

func (repo *TransactionRepository) GetReport(startDateStr string, endDateStr string) (*model.Report, error) {
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endDate := startDate.Add(24 * time.Hour)

	var err error

	// Parse start_date or default to today
	if startDateStr == "" {
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	} else {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {

		}
	}

	// Parse end_date or default to end of today
	if endDateStr == "" {
		now := time.Now()
		endDate = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	} else {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return nil, errors.New("Invalid end_date format. Use YYYY-MM-DD")
		}
		// Set to end of day
		endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, endDate.Location())
	}

	// Validate: start_date must be before end_date
	if startDate.After(endDate) {
		return nil, errors.New("start_date must be before end_date")
	}

	var transactions []model.Transaction
	var transactionID []int

	err = repo.builder.
		From("transactions").
		Where(goqu.I("created_at").Gte(startDate)).
		Where(goqu.I("created_at").Lt(endDate)).
		ScanStructs(&transactions)

	for _, transaction := range transactions {
		transactionID = append(transactionID, transaction.ID)
	}

	if err != nil {
		return nil, err
	}

	if len(transactionID) == 0 {
		return nil, nil
	}

	var report model.Report
	var productTerlaris model.ProductTerlaris

	_, err = repo.builder.
		From("transaction_details").
		Select(
			goqu.SUM("subtotal").As("total_revenue"),
			goqu.COUNT("id").As("total_transaksi"),
		).
		Where(goqu.C("transaction_id").Eq(transactionID)).
		ScanStruct(&report)

	_, err = repo.builder.
		From(goqu.T("transaction_details").As("td")).
		InnerJoin(
			goqu.T("products").As("p"),
			goqu.On(goqu.I("td.product_id").Eq(goqu.I("p.id"))),
		).
		Select(
			goqu.I("p.name").As("product_name"),
			goqu.SUM("td.quantity").As("total_qty"),
		).
		GroupBy("td.product_id", "p.name").
		Order(goqu.I("total_qty").Desc()).
		Limit(1).
		ScanStruct(&productTerlaris)

	report.ProductTerlaris = productTerlaris
	return &report, err

}
