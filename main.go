package main

import (
	"log"
	"net/http"

	"category-crud/config"
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
	config, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	log.Default().Println(config)
	store := NewCategoryStore()
	r := SetupRoutes(store)

	port := config.Server.Port
	log.Println("Server starting on :" + port)
	log.Println("Swagger documentation available at http://localhost:8080/swagger/index.html")
	log.Fatal(http.ListenAndServe(":"+port, r))
}
