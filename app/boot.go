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
	db, err := db.Configure(*config)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	handlerGroup := &handler.HandlerGroup{
		Product:  setupProduct(db),
		Category: setupCategory(db),
	}
	r := route.Configure(handlerGroup)

	port := config.Server.Port
	fmt.Println("Server starting on :" + port)
	fmt.Println("Swagger documentation available at http://localhost:8080/swagger/index.html")
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func setupProduct(db *sql.DB) *handler.ProductHandler {
	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	return productHandler
}

func setupCategory(db *sql.DB) *handler.CategoryHandler {
	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	return categoryHandler
}
