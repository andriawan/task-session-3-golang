package main

import (
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// SetupRoutes configures all API routes
func SetupRoutes(store *CategoryStore) *mux.Router {
	r := mux.NewRouter()

	// Root route - redirect to Swagger
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
	}).Methods("GET")

	// Category endpoints
	r.HandleFunc("/categories", store.CreateCategory).Methods("POST")
	r.HandleFunc("/categories", store.ListCategories).Methods("GET")
	r.HandleFunc("/categories/{id}", store.GetCategory).Methods("GET")
	r.HandleFunc("/categories/{id}", store.UpdateCategory).Methods("PUT")
	r.HandleFunc("/categories/{id}", store.DeleteCategory).Methods("DELETE")

	// Swagger documentation
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return r
}
