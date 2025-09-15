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
	mux.HandleFunc("GET /products/{id}", handlers.GetProductHandler)
	mux.HandleFunc("GET /collections", handlers.GetCollectionsHandler)
	mux.HandleFunc("GET /collections/{id}", handlers.GetCollectionHandler)
	mux.HandleFunc("GET /promotions", handlers.GetPromotionsHandler)

	mux.HandleFunc("POST /products", handlers.AddProductsHandler)
	mux.HandleFunc("POST /collections", handlers.AddCollectionsHandler)

	mux.HandleFunc("PUT /products/{id}", handlers.UpdateProductHandler)
	mux.HandleFunc("PUT /collections/{id}", handlers.UpdateCollectionHandler)

	mux.HandleFunc("PATCH /products/{id}", handlers.PatchProductHandler)
	mux.HandleFunc("PATCH /products", handlers.PatchProductsHandler)
	mux.HandleFunc("PATCH /collections/{id}", handlers.PactchCollectionHandler)
	mux.HandleFunc("PATCH /collections", handlers.PatchCollectionsHandler)

	mux.HandleFunc("DELETE /products/{id}", handlers.DeleteProductHandler)
	mux.HandleFunc("DELETE /products", handlers.DeleteProductsHandler)
	mux.HandleFunc("DELETE /collections/{id}", handlers.DeleteCollectionHandler)
	mux.HandleFunc("DELETE /collections", handlers.DeleteCollectionsHandler)
	

	return mux
}
