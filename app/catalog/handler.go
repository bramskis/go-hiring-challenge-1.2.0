package catalog

import (
	"encoding/json"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/models"
)

// having the response wrapper may not be necessary, but I will leave it to ensure that any consumers of this api do not break

type Response struct {
	Products []ProductCondensed `json:"products"`
}

// similar here, depending on the use case, we may not need to have a condensed version, but only sending the necessary information is not a bad idea

type ProductCondensed struct {
	Code     string  `json:"code"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

// do we need the dependency injection?

type CatalogHandler struct {
	repo *models.ProductsRepository
}

func NewCatalogHandler(r *models.ProductsRepository) *CatalogHandler {
	return &CatalogHandler{
		repo: r,
	}
}

func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	dbProducts, err := h.repo.GetAllProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	categoryIDSet := make(map[uint]any)
	for _, product := range dbProducts {
		if _, ok := categoryIDSet[product.CategoryID]; !ok {
			categoryIDSet[product.CategoryID] = nil
		}
	}

	categoryIDSlice := make([]uint, len(categoryIDSet))
	for id := range categoryIDSet {
		categoryIDSlice = append(categoryIDSlice, id)
	}

	categories, err := h.repo.GetCategoriesByID(categoryIDSlice)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Map response
	products := make([]ProductCondensed, len(dbProducts))
	for _, p := range dbProducts {
		var category string
		if _, ok := categories[p.CategoryID]; ok {
			category = categories[p.CategoryID].Name
		}
		products = append(products, ProductCondensed{
			Code:     p.Code,
			Price:    p.Price.InexactFloat64(),
			Category: category,
		})
	}

	// Return the products as a JSON response
	w.Header().Set("Content-Type", "application/json")

	response := Response{
		Products: products,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
