package handlers

import (
	"fmt"
	"net/http"
)

func GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Getting the product values")
}
