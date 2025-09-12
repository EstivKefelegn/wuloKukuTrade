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

func GetCollectionsHandler(w http.ResponseWriter, r *http.Request) {
	var Collections []models.Collection
	collections, totalCollection, err := sqlconnect.GetCollectionsDBHandler(Collections, r)
	if err != nil {
		return
	}

	response := struct {
		Status string              `json:"status"`
		Count  int                 `json:"count"`
		Data   []models.Collection `json:"data"`
	}{
		Status: "Succces",
		Count:  totalCollection,
		Data:   collections,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func AddCollectionsHandler(w http.ResponseWriter, r *http.Request) {
	var newCollections []models.Collection
	var rawCollections []map[string]interface{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error: invalid request body", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &rawCollections)
	if err != nil {
		http.Error(w, "Error: invalid request body", http.StatusBadRequest)
		return
	}

	fields := GetFieldNames(models.Collection{})
	allowedFields := make(map[string]struct{})

	for _, fields := range fields {
		allowedFields[fields] = struct{}{}
	}

	for _, collection := range rawCollections {
		for key := range collection {
			_, ok := allowedFields[key]
			if !ok {
				utils.ErrorHandler(err, "Unacceptable filed in the request")
				return
			}
		}
	}

	err = json.Unmarshal(body, &newCollections)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON format: %v", err), http.StatusBadRequest)
		return
	}

	for _, collection := range newCollections {
		err := CheckEmptyields(collection)
		if err != nil {
			return
		}
	}

	addCollection, err := sqlconnect.AddCollectionDBHandler(newCollections)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string
		Count  int
		Data   []models.Collection
	}{
		Status: "Success",
		Count:  len(addCollection),
		Data:   addCollection,
	}

	json.NewEncoder(w).Encode(response)

}
