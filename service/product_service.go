package service

import (
	"category-crud/model"
	"category-crud/model/dto"
	"category-crud/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAll(filter *dto.ProductFilterRequest) ([]model.Product, error) {
	return s.repo.GetAll(filter)
}

func (s *ProductService) Create(data *dto.ProductRequest) error {
	return s.repo.Create(data)
}

func (s *ProductService) GetByID(id int) (*model.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) Update(product *dto.ProductRequest) error {
	return s.repo.Update(product)
}

func (s *ProductService) Delete(id int) error {
	return s.repo.Delete(id)
}
