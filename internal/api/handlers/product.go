package handlers

import (
	"chickenTrade/API/internal/models"
	"chickenTrade/API/internal/repository/sqlconnect"
	"chickenTrade/API/pkg/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

func GetProductHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	product, err := sqlconnect.GetOneProductDBHandler(id)
	fmt.Println("Predouct: ", product)
	if err != nil {
		utils.ErrorHandler(err, "There is no product with this id")
		return
	}
	response := struct {
		Status string
		Data   models.Product
	}{
		Status: "Success",
		Data:   product,
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

func UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var updatedProduct models.Product
	err := json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		utils.ErrorHandler(err, "Couldnt convert the incomming request")
		return
	}

	updatedProductDB, err := sqlconnect.UpdateProductDBHandler(id, updatedProduct)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedProductDB)
}

func PatchProductHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var updatedProduct map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		utils.ErrorHandler(err, "couldn't decode the comming data")
		return
	}

	updatedProductDB, err := sqlconnect.PatchProductDBHandler(id, updatedProduct)
	if err != nil {
		log.Println(err)
		return
	}

	json.NewEncoder(w).Encode(updatedProductDB)

}

func PatchProductsHandler(w http.ResponseWriter, r *http.Request) {
	var updates []map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		utils.ErrorHandler(err, "couldn't parese the incomming requests")
		return
	}

	updatedProducts, err := sqlconnect.PatchProductsDBHandler(w, updates)
	if err != nil {
		log.Println(err)
		return
	}

	json.NewEncoder(w).Encode(updatedProducts)

}

func DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	id, err := sqlconnect.DeleteProductDBHandler(w, id)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status string `json:"status"`
		ID     string `json:"id"`
	}{
		Status: "Product Successfully deleted",
		ID:     id,
	}
	json.NewEncoder(w).Encode(response)

}

func DeleteProductsHandler(w http.ResponseWriter, r *http.Request) {
	var ids []string
	err := json.NewDecoder(r.Body).Decode(&ids)
	fmt.Println("ID's", ids)
	if err != nil {
		return
	}

	deletedIDs, err := sqlconnect.DeleteProductsDBHandler(w, ids)
	if err != nil {
		log.Println("Error: can't delete the products")
		return
	}

	json.NewEncoder(w).Encode(deletedIDs)
}
