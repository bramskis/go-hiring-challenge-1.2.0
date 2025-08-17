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

func (r *ProductsRepository) getCategoriesMappingByID(ids []uint) (map[uint]string, error) {
	var categories map[uint]Category
	if err := r.db.Find(&categories, ids).Error; err != nil {
		return nil, err
	}
	categoriesInfo := make(map[uint]string, len(categories))
	for _, category := range categories {
		if _, ok := categoriesInfo[category.ID]; !ok {
			categoriesInfo[category.ID] = category.Name
		}
	}
	return categoriesInfo, nil
}

type Catalog struct {
	Products        []Product
	CategoryDetails map[uint]string
	Total           int64
}

func (r *ProductsRepository) GetCatalogWithConditions(conditions GetCatalogConditions) (*Catalog, error) {
	catalog := Catalog{}

	query := r.db.Model(&Product{}).Joins("inner join on products.category_id = categories.id")

	if conditions.PriceLessThan != nil {
		query = query.Where("products.price < ?", *conditions.PriceLessThan)
	}
	if conditions.CategoryFilter != nil {
		query = query.Where("categories.name IN ?", *conditions.CategoryFilter)
	}

	query.Count(&catalog.Total)

	err := query.Offset(conditions.Offset).Limit(conditions.Limit).Find(&catalog.Products).Error
	if err != nil {
		return nil, err
	}

	categoryIDs := []uint{}
	categoryIDSet := make(map[uint]bool)
	for _, product := range catalog.Products {
		if _, ok := categoryIDSet[product.CategoryID]; !ok {
			categoryIDs = append(categoryIDs, product.CategoryID)
			categoryIDSet[product.CategoryID] = true
		}
	}

	catalog.CategoryDetails, err = r.getCategoriesMappingByID(categoryIDs)
	if err != nil {
		return nil, err
	}

	return &catalog, nil
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
