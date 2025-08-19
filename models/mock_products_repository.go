package models

import (
	"github.com/stretchr/testify/mock"
)

type MockProductsRepository struct {
	mock.Mock
}

func NewMockProductsRepository() MockProductsRepository {
	return MockProductsRepository{}
}

func (r *MockProductsRepository) GetProductByCode(code string) (*Product, error) {
	args := r.Called(code)
	if _, ok := args.Get(1).(error); ok {
		return nil, args.Get(1).(error)
	}
	return args.Get(0).(*Product), nil
}

func (r *MockProductsRepository) GetCatalogWithConditions(conditions GetCatalogConditions) (*Catalog, error) {
	args := r.Called(conditions)
	if _, ok := args.Get(1).(error); ok {
		return nil, args.Get(1).(error)
	}
	return args.Get(0).(*Catalog), nil
}

func (r *MockProductsRepository) GetAllCategories() ([]Category, error) {
	args := r.Called()
	if _, ok := args.Get(1).(error); ok {
		return nil, args.Get(1).(error)
	}
	return args.Get(0).([]Category), nil
}

func (r *MockProductsRepository) GetCategoryByID(id uint) (*Category, error) {
	args := r.Called(id)
	if _, ok := args.Get(1).(error); ok {
		return nil, args.Get(1).(error)
	}
	return args.Get(0).(*Category), nil
}

func (r *MockProductsRepository) CreateCategory(request CreateCategoryRequest) (*Category, error) {
	args := r.Called(request)
	if _, ok := args.Get(1).(error); ok {
		return nil, args.Get(1).(error)
	}
	return args.Get(0).(*Category), nil
}
