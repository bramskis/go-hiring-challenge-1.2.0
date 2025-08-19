package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type GetCatalogConditions struct {
	Offset         int
	Limit          int
	CategoryFilter []string
	PriceLessThan  decimal.Decimal
}

type ProductsRepository interface {
	GetProductByCode(code string) (*Product, error)
	GetCatalogWithConditions(conditions GetCatalogConditions) (*Catalog, error)
	GetAllCategories() ([]Category, error)
	GetCategoryByID(id uint) (*Category, error)
	CreateCategory(request CreateCategoryRequest) (*Category, error)
}

type ProductsRepositoryWithDB struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) ProductsRepository {
	return &ProductsRepositoryWithDB{
		db: db,
	}
}

func (r *ProductsRepositoryWithDB) GetProductByCode(code string) (*Product, error) {
	var product Product
	if err := r.db.Preload("Variants").First(&product, "code = ?", code).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductsRepositoryWithDB) getCategoriesMappingByID(ids []uint) (map[uint]string, error) {
	var categories []Category
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

func (r *ProductsRepositoryWithDB) GetCatalogWithConditions(conditions GetCatalogConditions) (*Catalog, error) {
	catalog := Catalog{}

	query := r.db.Model(&Product{})

	if conditions.PriceLessThan.GreaterThan(decimal.Zero) {
		query = query.Where("products.price < ?", conditions.PriceLessThan.InexactFloat64())
	}
	if len(conditions.CategoryFilter) > 0 {
		query = query.Joins("inner join categories on products.category_id = categories.id").Where("categories.name IN ?", conditions.CategoryFilter)
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

func (r *ProductsRepositoryWithDB) GetAllCategories() ([]Category, error) {
	var categories []Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *ProductsRepositoryWithDB) GetCategoryByID(id uint) (*Category, error) {
	var category Category
	if err := r.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

type CreateCategoryRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func (r *ProductsRepositoryWithDB) CreateCategory(request CreateCategoryRequest) (*Category, error) {
	category := Category{
		Code: request.Code,
		Name: request.Name,
	}
	if err := r.db.Create(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}
