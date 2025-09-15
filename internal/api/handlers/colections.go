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

func GetCollectionHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	collection, err := sqlconnect.GetCollectionDBHandler(id)
	if err != nil {
		utils.ErrorHandler(err, "There is no product found")
		return
	}

	response := struct {
		Status string
		Data   models.Collection
	}{
		Status: "Success",
		Data:   collection,
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

func UpdateCollectionHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var updatedCollection models.Collection

	err := json.NewDecoder(r.Body).Decode(&updatedCollection)
	if err != nil {
		log.Println("Couldnt decode the incomming request")
		return
	}
	fmt.Println("ID:", id)
	updatedCollectionDB, err := sqlconnect.UpdateCollectionDBHandler(id, updatedCollection)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedCollectionDB)
}

func PactchCollectionHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var requestBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&requestBody)

	if err != nil {
		// http.Error(w, "Couldnt padrse the incomming data", http.StatusBadRequest)
		utils.ErrorHandler(err, "couldn't padrse the in comming data")
		return
	}

	collection, err := sqlconnect.PatchCollectionDBHandler(id, requestBody)
	if err != nil {
		// http.Error(w, "Couldnt update the data", http.StatusInternalServerError)
		utils.ErrorHandler(err, "couldn't padrse the in comming data")
		return
	}

	response := struct {
		Status string
		Data   models.Collection
	}{
		Status: "Success",
		Data:   collection,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func PatchCollectionsHandler(w http.ResponseWriter, r *http.Request) {
	var updates []map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updates)

	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		log.Fatal("Couldnt parse the incomming data", http.StatusBadRequest)
		return
	}
	updatedCollections, err := sqlconnect.PatchCollectionsDBHandler(w, updates)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Fatal("Couldnt update the  data", http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string
		Data   []models.Collection
	}{
		Status: "success",
		Data: updatedCollections,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}




func DeleteCollectionHandler(w http.ResponseWriter, r *http.Request)  {
	id := r.PathValue("id")

	DeletedID, err := sqlconnect.DeleteCollectionDBHandler(w, id)
	if err != nil {
		http.Error(w, "Couldn't delete the requested value", http.StatusInternalServerError)
		log.Fatal("Couldnt delted the requested value", http.StatusInternalServerError)
		return
	}

	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(DeletedID)
}

func DeleteCollectionsHandler(w http.ResponseWriter, r *http.Request)  {
	var deltedIDs []string

	json.NewDecoder(r.Body).Decode(&deltedIDs)

	deletedIds, err := sqlconnect.DeleteCollectionsDBHandler(w, deltedIDs)
	if err != nil {
		http.Error(w, "couldn't delete the reuested values", http.StatusInternalServerError)
		log.Fatalln("couldn't delete the requested values")
		return
	}

	response := struct{
		Status string
		Data []string
	}{
		Status: "Success",
		Data: deletedIds,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}