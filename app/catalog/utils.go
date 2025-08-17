package catalog

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
)

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

func getCategoryFilter(r *http.Request) *[]string {
	var categories []string
	categoriesStr := r.URL.Query().Get("category")
	if categoriesStr == "" {
		return nil
	}

	if strings.Contains(categoriesStr, ",") {
		categories = strings.Split(categoriesStr, ",")
	} else {
		categories = []string{categoriesStr}
	}
	return &categories
}

func getPriceLessThanFilter(r *http.Request) *decimal.Decimal {
	priceLessThanStr := r.URL.Query().Get("priceLessThan")
	if priceLessThanStr == "" {
		return nil
	}
	priceLessThan, err := decimal.NewFromString(priceLessThanStr)
	if err != nil {
		return nil
	}

	return &priceLessThan
}
