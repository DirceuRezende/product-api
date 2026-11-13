package server

import (
	"net/http"

	"products-api/internal/product"
)

func NewRouter(products *product.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Hello World"))
	})

	mux.HandleFunc("GET /products", products.List)
	mux.HandleFunc("GET /products/{id}", products.Get)
	mux.HandleFunc("POST /products", products.Create)
	mux.HandleFunc("PUT /products/{id}", products.Update)
	mux.HandleFunc("DELETE /products/{id}", products.Delete)

	return mux
}
