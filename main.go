package main

import (
	"log"
	"net/http"

	_ "category-crud/docs"
)

// @title Category CRUD API
// @version 1.0
// @description API for managing categories with CRUD operations
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
func main() {
	store := NewCategoryStore()
	r := SetupRoutes(store)

	log.Println("Server starting on :8080")
	log.Println("Swagger documentation available at http://localhost:8080/swagger/index.html")
	log.Fatal(http.ListenAndServe(":8080", r))
}
