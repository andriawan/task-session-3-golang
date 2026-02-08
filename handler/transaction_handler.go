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

// Report Transaction Today godoc
// @Summary Report Transaction Today
// @Description Report Transaction Today
// @Tags transaction
// @Produce json
// @Success 200 {array} model.Report "Report"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/report/hari-ini [get]
func (h *TransactionHandler) GetReportToday(w http.ResponseWriter, r *http.Request) {
	report, err := h.service.GetReport("", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// Report Transaction Based on Date godoc
// @Summary Report Transaction Based on Date
// @Description Report Transaction Based on Date
// @Tags transaction
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)" example(2026-01-01)
// @Param end_date query string false "End date (YYYY-MM-DD)" example(2026-02-01)
// @Success 200 {array} model.Report "Report"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/report [get]
func (h *TransactionHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	report, err := h.service.GetReport(startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}
