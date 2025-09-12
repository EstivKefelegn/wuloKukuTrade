package router

import (
	"fmt"
	"net/http"
)

func CollectionsRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /collections", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Getting the collections")
	})

	return mux
}
