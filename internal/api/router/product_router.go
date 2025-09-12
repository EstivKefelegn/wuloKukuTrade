package router

import (
	"fmt"
	"net/http"
)

func ProductRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /products", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Getting the products of the KUKU Trade")
	})

	return mux
}
