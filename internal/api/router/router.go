package router

import (
	"chickenTrade/API/internal/api/handlers"
	"fmt"
	"net/http"
)

func testRouter(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "KUKU trade main router stars")
}

func MainRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", testRouter)

	mux.HandleFunc("GET /products", handlers.GetProductsHandler)
	mux.HandleFunc("GET /collections", handlers.GetCollectionsHandler)
	mux.HandleFunc("GET /promotions", handlers.GetPromotionsHandler)	

	return mux
}
