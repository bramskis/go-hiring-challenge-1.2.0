package catalog

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setMockProductsRepository(repo *models.MockProductsRepository) {
	product1 := models.Product{
		ID:         1,
		Code:       "testCode1",
		Price:      decimal.NewFromFloat(12.34),
		CategoryID: 1,
		Variants: []models.Variant{
			{
				ID:    1,
				Name:  "testVariant1",
				SKU:   "testSKU1",
				Price: decimal.Zero,
			},
			{
				ID:    2,
				Name:  "testVariant2",
				SKU:   "testSKU2",
				Price: decimal.NewFromFloat(23.45),
			},
		},
	}
	repo.On("GetProductByCode", mock.Anything).Return(&product1, nil)

	product2 := models.Product{
		ID:         1,
		Code:       "testCode2",
		Price:      decimal.NewFromFloat(23.45),
		CategoryID: 2,
	}
	catalog := models.Catalog{
		Products: []models.Product{
			product1,
			product2,
		},
		CategoryDetails: map[uint]string{
			1: "testCategory1",
			2: "testCategory2",
		},
		Total: 2,
	}
	repo.On("GetCatalogWithConditions", mock.Anything).Return(&catalog, nil)

	category1 := models.Category{
		ID:   1,
		Code: "testCode1",
		Name: "testName1",
	}
	category2 := models.Category{
		ID:   2,
		Code: "testCode2",
		Name: "testName2",
	}
	categories := []models.Category{
		category1,
		category2,
	}
	repo.On("GetAllCategories").Return(categories, nil)

	repo.On("GetCategoryByID", mock.Anything).Return(&category1, nil)

	repo.On("CreateCategory", mock.Anything).Return(&category2, nil)
}

func getMockProductsRepositoryWithResponses() *models.MockProductsRepository {
	repo := models.NewMockProductsRepository()
	setMockProductsRepository(&repo)
	return &repo
}

func getMockProductsRepositoryNoResponses() *models.MockProductsRepository {
	repo := models.NewMockProductsRepository()
	return &repo
}

func TestHandleGetCatalog_HappyPath(t *testing.T) {
	t.Run("Successfully returns catalog info", func(t *testing.T) {
		// Setup
		recorder := httptest.NewRecorder()
		repo := getMockProductsRepositoryWithResponses()
		handler := CatalogHandler{repo: repo}
		req, err := http.NewRequest(http.MethodGet, "http://example.com/catalog", nil)
		assert.Nil(t, err, "unable to create http request")
		handler.HandleGetCatalog(recorder, req)

		repo.AssertCalled(t, "GetCatalogWithConditions", models.GetCatalogConditions{
			Offset:         0,
			Limit:          10,
			CategoryFilter: []string{},
			PriceLessThan:  decimal.Zero,
		})
		assert.Equal(t, http.StatusOK, recorder.Code, "Expected status code 200 OK")
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"), "Expected Content-Type to be application/json")

		expected, err := json.Marshal(GetCatalogResponse{
			Products: []ProductCondensed{
				{
					Code:     "testCode1",
					Price:    12.34,
					Category: "testCategory1",
				},
				{
					Code:     "testCode2",
					Price:    23.45,
					Category: "testCategory2",
				},
			},
			TotalProducts: 2,
		})
		assert.Nil(t, err)
		assert.JSONEq(t, string(expected), recorder.Body.String(), "Response body does not match expected")
	})
	t.Run("Successfully retrieves all query params", func(t *testing.T) {
		// Setup
		recorder := httptest.NewRecorder()
		repo := getMockProductsRepositoryWithResponses()
		handler := CatalogHandler{repo: repo}
		req, err := http.NewRequest(
			http.MethodGet,
			"http://example.com/catalog?offset=2&limit=3&category=Accessories,Shoes&priceLessThan=23.45",
			nil,
		)
		assert.Nil(t, err, "unable to create http request")
		handler.HandleGetCatalog(recorder, req)

		repo.AssertCalled(t, "GetCatalogWithConditions", models.GetCatalogConditions{
			Offset:         2,
			Limit:          3,
			CategoryFilter: []string{"Accessories", "Shoes"},
			PriceLessThan:  decimal.NewFromFloat(23.45),
		})
		assert.Equal(t, http.StatusOK, recorder.Code, "Expected status code 200 OK")
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"), "Expected Content-Type to be application/json")

		expected, err := json.Marshal(GetCatalogResponse{
			Products: []ProductCondensed{
				{
					Code:     "testCode1",
					Price:    12.34,
					Category: "testCategory1",
				},
				{
					Code:     "testCode2",
					Price:    23.45,
					Category: "testCategory2",
				},
			},
			TotalProducts: 2,
		})
		assert.Nil(t, err)
		assert.JSONEq(t, string(expected), recorder.Body.String(), "Response body does not match expected")
	})
}

func TestHandleGetCatalog_ErrorPath(t *testing.T) {
	t.Run("Handles not found errors", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		repo := getMockProductsRepositoryNoResponses()
		(*repo).On("GetCatalogWithConditions", mock.Anything).Return(nil, errors.New("record not found"))

		handler := CatalogHandler{repo: repo}
		req, err := http.NewRequest(http.MethodGet, "http://example.com/catalog", nil)
		assert.Nil(t, err, "unable to create http request")
		handler.HandleGetCatalog(recorder, req)

		repo.AssertCalled(t, "GetCatalogWithConditions", models.GetCatalogConditions{
			Offset:         0,
			Limit:          10,
			CategoryFilter: []string{},
			PriceLessThan:  decimal.Zero,
		})
		assert.Equal(t, http.StatusBadRequest, recorder.Code, "Expected status code 400")
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"), "Expected Content-Type to be application/json")

		expected, err := json.Marshal(api.ErrorJSON{Error: "record not found"})
		assert.Nil(t, err)
		assert.JSONEq(t, string(expected), recorder.Body.String(), "Response body does not match expected")
	})
	t.Run("Handles other errors", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		repo := getMockProductsRepositoryNoResponses()
		repo.On("GetCatalogWithConditions", mock.Anything).Return(nil, errors.New("test error"))

		handler := CatalogHandler{repo: repo}
		req, err := http.NewRequest(http.MethodGet, "http://example.com/catalog", nil)
		assert.Nil(t, err, "unable to create http request")
		handler.HandleGetCatalog(recorder, req)

		repo.AssertCalled(t, "GetCatalogWithConditions", models.GetCatalogConditions{
			Offset:         0,
			Limit:          10,
			CategoryFilter: []string{},
			PriceLessThan:  decimal.Zero,
		})
		assert.Equal(t, http.StatusInternalServerError, recorder.Code, "Expected status code 500")
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"), "Expected Content-Type to be application/json")

		expected, err := json.Marshal(api.ErrorJSON{Error: "test error"})
		assert.Nil(t, err)
		assert.JSONEq(t, string(expected), recorder.Body.String(), "Response body does not match expected")
	})
}

func TestHandleGetProductByCode_HappyPath(t *testing.T) {
	t.Run("Successfully returns product info including variants", func(t *testing.T) {
		// Setup
		recorder := httptest.NewRecorder()
		repo := getMockProductsRepositoryWithResponses()
		handler := CatalogHandler{repo: repo}
		req, err := http.NewRequest(http.MethodGet, "http://example.com/catalog/testCode1", nil)
		assert.Nil(t, err, "unable to create http request")
		handler.HandleGetProductByCode(recorder, req)

		repo.AssertCalled(t, "GetProductByCode", "testCode1")
		assert.Equal(t, http.StatusOK, recorder.Code, "Expected status code 200 OK")
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"), "Expected Content-Type to be application/json")

		expected, err := json.Marshal(GetProductByCodeResponse{
			Product: ProductCondensed{
				Code:     "testCode1",
				Price:    12.34,
				Category: "testName1",
			},
			Variants: []VariantCondensed{
				{
					Name:  "testVariant1",
					SKU:   "testSKU1",
					Price: 12.34,
				},
				{
					Name:  "testVariant2",
					SKU:   "testSKU2",
					Price: 23.45,
				},
			},
		})
		assert.Nil(t, err)
		assert.JSONEq(t, string(expected), recorder.Body.String(), "Response body does not match expected")
	})
}

func TestHandleGetProductByCode_ErrorPath(t *testing.T) {
	t.Run("Handles not found errors", func(t *testing.T) {
		// Setup
		recorder := httptest.NewRecorder()
		repo := getMockProductsRepositoryNoResponses()
		repo.On("GetProductByCode", mock.Anything).Return(nil, errors.New("record not found"))
		handler := CatalogHandler{repo: repo}
		req, err := http.NewRequest(http.MethodGet, "http://example.com/catalog/testCode1", nil)
		assert.Nil(t, err, "unable to create http request")
		handler.HandleGetProductByCode(recorder, req)

		repo.AssertCalled(t, "GetProductByCode", "testCode1")
		assert.Equal(t, http.StatusBadRequest, recorder.Code, "Expected status code 400")
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"), "Expected Content-Type to be application/json")

		expected, err := json.Marshal(api.ErrorJSON{Error: "record not found"})
		assert.Nil(t, err)
		assert.JSONEq(t, string(expected), recorder.Body.String(), "Response body does not match expected")
	})
	t.Run("Handles other errors", func(t *testing.T) {
		// Setup
		recorder := httptest.NewRecorder()
		repo := getMockProductsRepositoryNoResponses()
		repo.On("GetProductByCode", mock.Anything).Return(nil, errors.New("test error"))
		handler := CatalogHandler{repo: repo}
		req, err := http.NewRequest(http.MethodGet, "http://example.com/catalog/testCode1", nil)
		assert.Nil(t, err, "unable to create http request")
		handler.HandleGetProductByCode(recorder, req)

		repo.AssertCalled(t, "GetProductByCode", "testCode1")
		assert.Equal(t, http.StatusInternalServerError, recorder.Code, "Expected status code 500")
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"), "Expected Content-Type to be application/json")

		expected, err := json.Marshal(api.ErrorJSON{Error: "test error"})
		assert.Nil(t, err)
		assert.JSONEq(t, string(expected), recorder.Body.String(), "Response body does not match expected")
	})
}

func TestHandleGetCategories_HappyPath(t *testing.T) {
	t.Run("Successfully returns categories", func(t *testing.T) {
		// Setup
		recorder := httptest.NewRecorder()
		repo := getMockProductsRepositoryWithResponses()
		handler := CatalogHandler{repo: repo}
		req, err := http.NewRequest(http.MethodGet, "http://example.com/categories", nil)
		assert.Nil(t, err, "unable to create http request")
		handler.HandleGetCategories(recorder, req)

		repo.AssertCalled(t, "GetAllCategories")
		assert.Equal(t, http.StatusOK, recorder.Code, "Expected status code 200 OK")
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"), "Expected Content-Type to be application/json")

		expected, err := json.Marshal(GetCategoriesResponse{
			Categories: []CategoryCondensed{
				{
					ID:   1,
					Code: "testCode1",
					Name: "testName1",
				},
				{
					ID:   2,
					Code: "testCode2",
					Name: "testName2",
				},
			},
		})
		assert.Nil(t, err)
		assert.JSONEq(t, string(expected), recorder.Body.String(), "Response body does not match expected")
	})
}

func TestHandleGetCategories_ErrorPath(t *testing.T) {
	t.Run("Handles errors", func(t *testing.T) {
		// Setup
		recorder := httptest.NewRecorder()
		repo := getMockProductsRepositoryNoResponses()
		repo.On("GetAllCategories", mock.Anything).Return(nil, errors.New("test error"))
		handler := CatalogHandler{repo: repo}
		req, err := http.NewRequest(http.MethodGet, "http://example.com/categories", nil)
		assert.Nil(t, err, "unable to create http request")
		handler.HandleGetCategories(recorder, req)

		repo.AssertCalled(t, "GetAllCategories")
		assert.Equal(t, http.StatusInternalServerError, recorder.Code, "Expected status code 500")
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"), "Expected Content-Type to be application/json")

		expected, err := json.Marshal(api.ErrorJSON{Error: "test error"})
		assert.Nil(t, err)
		assert.JSONEq(t, string(expected), recorder.Body.String(), "Response body does not match expected")
	})
}

func TestHandlePostCategory_HappyPath(t *testing.T) {
	t.Run("Successfully creates and returns categories", func(t *testing.T) {
		// Setup
		recorder := httptest.NewRecorder()
		repo := getMockProductsRepositoryWithResponses()
		handler := CatalogHandler{repo: repo}
		body, err := json.Marshal(models.CreateCategoryRequest{
			Code: "testCode2",
			Name: "testName2",
		})
		assert.Nil(t, err)
		req, err := http.NewRequest(
			http.MethodPost,
			"http://example.com/categories",
			strings.NewReader(string(body)),
		)
		assert.Nil(t, err, "unable to create http request")
		handler.HandlePostCategory(recorder, req)

		repo.AssertCalled(t, "CreateCategory", models.CreateCategoryRequest{
			Code: "testCode2",
			Name: "testName2",
		})
		assert.Equal(t, http.StatusOK, recorder.Code, "Expected status code 200 OK")
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"), "Expected Content-Type to be application/json")

		expected, err := json.Marshal(PostCategoryResponse{
			Category: CategoryCondensed{
				ID:   2,
				Code: "testCode2",
				Name: "testName2",
			},
		})
		assert.Nil(t, err)
		assert.JSONEq(t, string(expected), recorder.Body.String(), "Response body does not match expected")
	})
}

func TestHandlePostCategory_ErrorPath(t *testing.T) {
	t.Run("Handles bad input", func(t *testing.T) {
		// Setup
		recorder := httptest.NewRecorder()
		repo := getMockProductsRepositoryWithResponses()
		handler := CatalogHandler{repo: repo}
		body, err := json.Marshal(VariantCondensed{}) // any incorrect struct would work
		assert.Nil(t, err)
		req, err := http.NewRequest(
			http.MethodPost,
			"http://example.com/categories",
			strings.NewReader(string(body)),
		)
		assert.Nil(t, err, "unable to create http request")
		handler.HandlePostCategory(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code, "Expected status code 400")
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"), "Expected Content-Type to be application/json")

		expected, err := json.Marshal(api.ErrorJSON{Error: "Request does not match required fields"})
		assert.Nil(t, err)
		assert.JSONEq(t, string(expected), recorder.Body.String(), "Response body does not match expected")
	})
}
