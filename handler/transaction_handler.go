package handler

import (
	"category-crud/model"
	"category-crud/service"
	"encoding/json"
	"net/http"
)

type TransactionHandler struct {
	service *service.TransactionService
}

func NewTransactionHandler(service *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// Checkout godoc
// @Summary Checkout products
// @Description Checkout selected products
// @Tags transaction
// @Accept json
// @Produce json
// @Success 200 {array} model.Transaction "Transaction"
// @Param request body model.CheckoutRequest true "Checkout payload"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/checkout [post]
func (h *TransactionHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	var req model.CheckoutRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	transaction, err := h.service.Checkout(req.Items)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}
