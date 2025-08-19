package catalog

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mytheresa/go-hiring-challenge/models"
)

type CodeResponse struct {
	Product  ProductCondensed   `json:"products"`
	Variants []VariantCondensed `json:"variants"`
}

type ProductCondensed struct {
	Code     string  `json:"code"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

type VariantCondensed struct {
	Name  string  `json:"name"`
	SKU   string  `json:"sku"`
	Price float64 `json:"price"`
}

func (h *CatalogHandler) HandleGetCode(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/catalog/")

	product, err := h.repo.GetProductByCode(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	category, err := h.repo.GetCategoryByID(product.CategoryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	response := CodeResponse{
		Product: ProductCondensed{
			Code:     product.Code,
			Price:    product.Price.InexactFloat64(),
			Category: category.Name,
		},
		Variants: variantsCondensed,
	}

	// Return the products as a JSON response
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type CatalogResponse struct {
	Products      []ProductCondensed `json:"products"`
	TotalProducts int                `json:"totalProducts"`
}

type CatalogHandler struct {
	repo *models.ProductsRepository
}

func NewCatalogHandler(r *models.ProductsRepository) *CatalogHandler {
	return &CatalogHandler{
		repo: r,
	}
}

func (h *CatalogHandler) HandleGetCatalog(w http.ResponseWriter, r *http.Request) {
	conditions := createGetCatalogConditions(r)

	catalog, err := h.repo.GetCatalogWithConditions(conditions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	// Return the products as a JSON response
	w.Header().Set("Content-Type", "application/json")

	response := CatalogResponse{
		Products: productsCondensed,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
