package catalog

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
)

type CatalogHandler struct {
	repo *models.ProductsRepository
}

func NewCatalogHandler(r *models.ProductsRepository) *CatalogHandler {
	return &CatalogHandler{
		repo: r,
	}
}

type GetCatalogResponse struct {
	Products      []ProductCondensed `json:"products"`
	TotalProducts int                `json:"totalProducts"`
}

func (h *CatalogHandler) HandleGetCatalog(w http.ResponseWriter, r *http.Request) {
	conditions := createGetCatalogConditions(r)

	catalog, err := h.repo.GetCatalogWithConditions(conditions)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") { // bad code in URI
			api.ErrorResponse(w, http.StatusBadRequest, err.Error())
		} else {
			api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	productsCondensed := make([]ProductCondensed, len(catalog.Products))
	for i, product := range catalog.Products {
		categoryName := catalog.CategoryDetails[product.CategoryID]
		productsCondensed[i] = ProductCondensed{
			Code:     product.Code,
			Price:    product.Price.InexactFloat64(),
			Category: categoryName,
		}
	}

	api.OKResponse(w, GetCatalogResponse{
		Products: productsCondensed,
	})
}

type GetProductByCodeResponse struct {
	Product  ProductCondensed   `json:"products"`
	Variants []VariantCondensed `json:"variants"`
}

func (h *CatalogHandler) HandleGetProductByCode(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/catalog/")

	product, err := h.repo.GetProductByCode(code)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") { // bad code in URI
			api.ErrorResponse(w, http.StatusBadRequest, err.Error())
		} else {
			api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	category, err := h.repo.GetCategoryByID(product.CategoryID)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error()) // server error or product has non-existent category
		return
	}

	variantsCondensed := make([]VariantCondensed, len(product.Variants))
	for i, variant := range product.Variants {
		variantPrice := product.Price.InexactFloat64()
		if variant.Price.InexactFloat64() > 0 {
			variantPrice = variant.Price.InexactFloat64()
		}
		variantsCondensed[i] = VariantCondensed{
			Name:  variant.Name,
			SKU:   variant.SKU,
			Price: variantPrice,
		}
	}

	api.OKResponse(w, GetProductByCodeResponse{
		Product: ProductCondensed{
			Code:     product.Code,
			Price:    product.Price.InexactFloat64(),
			Category: category.Name,
		},
		Variants: variantsCondensed,
	})
}

type GetCategoriesResponse struct {
	Categories []CategoryCondensed `json:"categories"`
}

func (h *CatalogHandler) HandleGetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.repo.GetAllCategories()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	categoriesCondensed := make([]CategoryCondensed, len(categories))
	for i, category := range categories {
		categoriesCondensed[i] = CategoryCondensed{
			ID:   category.ID,
			Code: category.Code,
			Name: category.Name,
		}
	}

	api.OKResponse(w, GetCategoriesResponse{
		Categories: categoriesCondensed,
	})
}

type PostCategoryResponse struct {
	Category CategoryCondensed `json:"category"`
}

func (h *CatalogHandler) HandlePostCategory(w http.ResponseWriter, r *http.Request) {
	createCategoryRequest := models.CreateCategoryRequest{}
	err := json.NewDecoder(r.Body).Decode(&createCategoryRequest)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	category, err := h.repo.CreateCategory(createCategoryRequest)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.OKResponse(w, PostCategoryResponse{
		Category: CategoryCondensed{
			ID:   category.ID,
			Code: category.Code,
			Name: category.Name,
		},
	})
}
