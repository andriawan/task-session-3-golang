package service

import (
	"category-crud/model"
	"category-crud/repository"
)

type TransactionService struct {
	repo *repository.TransactionRepository
}

func NewTransactionService(repo *repository.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Checkout(items []model.CheckoutItem) (*model.Transaction, error) {
	return s.repo.CreateTransaction(items)
}
