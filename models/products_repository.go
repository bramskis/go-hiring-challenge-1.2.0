package models

import (
	"gorm.io/gorm"
)

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetAllProducts() ([]Product, error) {
	var products []Product
	if err := r.db.Preload("Variants").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductsRepository) GetAllCategoriesAndProducts() ([]Category, error) {
	var categories []Category
	if err := r.db.Preload("Products").Preload("Variants").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *ProductsRepository) GetCategoriesByID(ids []uint) (map[uint]Category, error) {
	var categories map[uint]Category
	if err := r.db.Find(&categories, ids).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
