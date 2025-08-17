package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type GetCatalogConditions struct {
	Offset         int
	Limit          int
	CategoryFilter *[]string
	PriceLessThan  *decimal.Decimal
}

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

func (r *ProductsRepository) GetCatalog() ([]Category, error) {
	var categories []Category
	if err := r.db.Preload("Products").Preload("Variants").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *ProductsRepository) GetCatalogWithConditions(conditions GetCatalogConditions) ([]Category, error) {
	var categories []Category
	query := r.db.Preload("Products").Preload("Variants")
	if conditions.CategoryFilter != nil {
		query = query.Where("categories.name IN ?", *conditions.CategoryFilter)
	}
	if conditions.PriceLessThan != nil {
		query = query.Where("products.price < ?", *conditions.PriceLessThan)
	}
	err := query.Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *ProductsRepository) GetCategories() ([]Category, error) {
	var categories []Category
	if err := r.db.Find(&categories).Error; err != nil {
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
