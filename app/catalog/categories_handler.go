package catalog

import (
	"encoding/json"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/models"
)

type CategoryCondensed struct {
	ID   uint   `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type GetCategoriesResponse struct {
	Categories []CategoryCondensed `json:"categories"`
}

func (h *CatalogHandler) HandleGetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.repo.GetAllCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	categoriesCondensed := make([]CategoryCondensed, len(categories))
	for _, category := range categories {
		categoriesCondensed = append(categoriesCondensed, CategoryCondensed{
			ID:   category.ID,
			Code: category.Code,
			Name: category.Name,
		})
	}

	response := GetCategoriesResponse{
		Categories: categoriesCondensed,
	}

	// Return the products as a JSON response
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type PostCategoryResponse struct {
	Category CategoryCondensed `json:"category"`
}

func (h *CatalogHandler) HandlePostCategory(w http.ResponseWriter, r *http.Request) {
	createCategoryRequest := models.CreateCategoryRequest{}
	err := json.NewDecoder(r.Body).Decode(&createCategoryRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	category, err := h.repo.CreateCategory(createCategoryRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := PostCategoryResponse{
		Category: CategoryCondensed{
			ID:   category.ID,
			Code: category.Code,
			Name: category.Name,
		},
	}

	// Return the products as a JSON response
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
