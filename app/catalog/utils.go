package catalog

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
)

// Condensed types containing only information to be returned from API

type CategoryCondensed struct {
	ID   uint   `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
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

func createGetCatalogConditions(r *http.Request) models.GetCatalogConditions {
	return models.GetCatalogConditions{
		Offset:         getOffset(r),
		Limit:          getLimit(r),
		CategoryFilter: getCategoryFilter(r),
		PriceLessThan:  getPriceLessThanFilter(r),
	}
}

func getOffset(r *http.Request) int {
	var offset int
	var err error
	offsetStr := r.URL.Query().Get("offset")

	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			offset = 0
		}
	}
	return offset
}

func getLimit(r *http.Request) int {
	limit := 10
	var err error
	limitStr := r.URL.Query().Get("limit")

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit > 100 || limit < 1 {
			limit = 10
		}
	}
	return limit
}

func getCategoryFilter(r *http.Request) []string {
	categories := []string{}
	categoriesStr := r.URL.Query().Get("category")
	if categoriesStr == "" {
		return categories
	}

	if strings.Contains(categoriesStr, ",") {
		categories = strings.Split(categoriesStr, ",")
	} else {
		categories = []string{categoriesStr}
	}
	return categories
}

func getPriceLessThanFilter(r *http.Request) decimal.Decimal {
	priceLessThanStr := r.URL.Query().Get("priceLessThan")
	if priceLessThanStr == "" {
		return decimal.Zero
	}
	priceLessThan, err := decimal.NewFromString(priceLessThanStr)
	if err != nil {
		return decimal.Zero
	}

	return priceLessThan
}
