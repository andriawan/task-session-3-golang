package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CategoryStore struct {
	mu         sync.RWMutex
	categories map[int]*Category
	nextID     int
}

func NewCategoryStore() *CategoryStore {
	return &CategoryStore{
		categories: make(map[int]*Category),
		nextID:     1,
	}
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new category with name and description
// @Tags categories
// @Accept json
// @Produce json
// @Param category body Category true "Category object"
// @Success 201 {object} Category
// @Failure 400 {string} string "Bad Request"
// @Router /categories [post]
func (s *CategoryStore) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	category.ID = s.nextID
	s.nextID++
	s.categories[category.ID] = &category
	s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

// ListCategories godoc
// @Summary List all categories
// @Description Get a list of all categories
// @Tags categories
// @Produce json
// @Success 200 {array} Category
// @Router /categories [get]
func (s *CategoryStore) ListCategories(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	categories := make([]*Category, 0, len(s.categories))
	for _, category := range s.categories {
		categories = append(categories, category)
	}
	s.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// GetCategory godoc
// @Summary Get a category by ID
// @Description Get a single category by its ID
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} Category
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Category not found"
// @Router /categories/{id} [get]
func (s *CategoryStore) GetCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	s.mu.RLock()
	category, exists := s.categories[id]
	s.mu.RUnlock()

	if !exists {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// UpdateCategory godoc
// @Summary Update a category
// @Description Update an existing category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body Category true "Category object"
// @Success 200 {object} Category
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Category not found"
// @Router /categories/{id} [put]
func (s *CategoryStore) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	var category Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	if _, exists := s.categories[id]; !exists {
		s.mu.Unlock()
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	category.ID = id
	s.categories[id] = &category
	s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// DeleteCategory godoc
// @Summary Delete a category
// @Description Delete a category by ID
// @Tags categories
// @Param id path int true "Category ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Category not found"
// @Router /categories/{id} [delete]
func (s *CategoryStore) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	if _, exists := s.categories[id]; !exists {
		s.mu.Unlock()
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	delete(s.categories, id)
	s.mu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}
