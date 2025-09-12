package handlers

import (
	"fmt"
	"net/http"
)

func GetPromotionsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Getting the Promotion values")
}
