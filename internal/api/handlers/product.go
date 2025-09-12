package handlers

import (
	"chickenTrade/API/internal/models"
	"chickenTrade/API/internal/repository/sqlconnect"
	"chickenTrade/API/pkg/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	var Products []models.Product
	products, totalProducts, err := sqlconnect.GetProductsDBHandler(Products, r)
	if err != nil {
		return
	}

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Product `json:"data"`
	}{
		Status: "Success",
		Count:  totalProducts,
		Data:   products,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func AddProductsHandler(w http.ResponseWriter, r *http.Request) {
	var newProducts []models.Product
	var rawProducts []map[string]interface{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error: reading request body", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &rawProducts)
	if err != nil {
		http.Error(w, "Error: invalid request body", http.StatusBadRequest)
		return
	}

	fields := GetFieldNames(models.Product{})
	allowedFields := make(map[string]struct{})

	for _, field := range fields {
		allowedFields[field] = struct{}{}
	}

	for _, product := range rawProducts {
		for key := range product {
			_, ok := allowedFields[key]
			fmt.Println("KEYS: ", key)
			if !ok {
				utils.ErrorHandler(err, "Unacceptable filed in the request")
				return
			}
		}
	}

	err = json.Unmarshal(body, &newProducts)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON format: %v", err), http.StatusBadRequest)
		return
	}

	for _, product := range newProducts {
		err := CheckEmptyields(product)
		if err != nil {
			return
		}
	}

	addProducts, err := sqlconnect.AddProductsDBHandler(newProducts)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string
		Count  int
		Data   []models.Product
	}{
		Status: "Success",
		Count:  len(addProducts),
		Data:   addProducts,
	}

	json.NewEncoder(w).Encode(response)

}
