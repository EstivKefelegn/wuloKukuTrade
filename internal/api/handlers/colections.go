package handlers

import (
	"fmt"
	"net/http"
)

func GetCollectionsHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "Getting the Collection values")
}