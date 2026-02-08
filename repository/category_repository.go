package repository

import (
	"category-crud/model"
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v9"
)

type CategoryRepository struct {
	db      *sql.DB
	builder *goqu.Database
}

func NewCategoryRepository(db *sql.DB, builder *goqu.Database) *CategoryRepository {
	return &CategoryRepository{
		db:      db,
		builder: builder,
	}
}

func (repo *CategoryRepository) GetAll() ([]model.Category, error) {
	var categories []model.Category
	err := repo.builder.From("categories").
		Select("id", "name", "description").
		ScanStructs(&categories)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (repo *CategoryRepository) Create(category *model.Category) error {
	_, err := repo.builder.Insert("categories").Rows(
		goqu.Record{
			"name":        category.Name,
			"description": category.Description,
		},
	).Returning("id").Executor().ScanStruct(category)
	return err
}

// GetByID - ambil kategori by ID
func (repo *CategoryRepository) GetByID(id int) (*model.Category, error) {
	var category model.Category
	result, err := repo.builder.From("categories").
		Select("id", "name", "description").
		Where(goqu.Ex{
			"id": id,
		}).ScanStruct(&category)

	if !result {
		return nil, errors.New("category tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (repo *CategoryRepository) Update(category *model.Category) error {
	result, err := repo.builder.Update("categories").Set(
		goqu.Record{
			"name":        category.Name,
			"description": category.Description,
		},
	).Where(goqu.Ex{
		"id": category.ID,
	}).Executor().Exec()

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("kategori tidak ditemukan")
	}

	return err
}

func (repo *CategoryRepository) Delete(id int) error {
	result, err := repo.builder.Delete("categories").Where(goqu.Ex{
		"id": id,
	}).Executor().Exec()
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("kategori tidak ditemukan")
	}

	return err
}
