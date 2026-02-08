package app

import (
	"category-crud/config"
	"category-crud/db"
	_ "category-crud/docs"
	"category-crud/handler"
	"category-crud/repository"
	"category-crud/route"
	"category-crud/service"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/doug-martin/goqu/v9"
)

// @title Category CRUD API
// @version 1.0
// @description API for managing categories with CRUD operations
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @BasePath /

func Start() {
	config, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	db, builder, err := db.Configure(*config)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	productHandler, _, productRepo := setupProduct(db, builder)
	handlerGroup := &handler.HandlerGroup{
		Product:     productHandler,
		Category:    setupCategory(db, builder),
		Transaction: setupTransaction(db, builder, productRepo),
	}
	r := route.Configure(handlerGroup)

	port := config.Server.Port
	fmt.Println("Server starting on :" + port)
	fmt.Println("Swagger documentation available at http://localhost:8080/swagger/index.html")
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func setupProduct(db *sql.DB, builder *goqu.Database) (*handler.ProductHandler, *service.ProductService, *repository.ProductRepository) {
	productRepo := repository.NewProductRepository(db, builder)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	return productHandler, productService, productRepo
}

func setupCategory(db *sql.DB, builder *goqu.Database) *handler.CategoryHandler {
	categoryRepo := repository.NewCategoryRepository(db, builder)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	return categoryHandler
}

func setupTransaction(db *sql.DB, builder *goqu.Database, productRepo *repository.ProductRepository) *handler.TransactionHandler {
	transactionRepo := repository.NewTransactionRepository(db, builder, productRepo)
	transactionService := service.NewTransactionService(transactionRepo)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	return transactionHandler
}
