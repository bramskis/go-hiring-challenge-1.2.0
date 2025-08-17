package catalog

import (
	"encoding/json"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/models"
)

type Response struct {
	Products      []ProductCondensed `json:"products"`
	TotalProducts int                `json:"totalProducts"`
}

type ProductCondensed struct {
	Code     string  `json:"code"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

type CatalogHandler struct {
	repo *models.ProductsRepository
}

func NewCatalogHandler(r *models.ProductsRepository) *CatalogHandler {
	return &CatalogHandler{
		repo: r,
	}
}

func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	conditions := createGetCatalogConditions(r)

	catalog, err := h.repo.GetCatalogWithConditions(conditions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	productsCondensed := make([]ProductCondensed, len(catalog.Products))
	for _, product := range catalog.Products {
		categoryName := catalog.CategoryDetails[product.CategoryID]
		productsCondensed = append(productsCondensed, ProductCondensed{
			Code:     product.Code,
			Price:    product.Price.InexactFloat64(),
			Category: categoryName,
		})
	}

	// Return the products as a JSON response
	w.Header().Set("Content-Type", "application/json")

	response := Response{
		Products: productsCondensed,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
