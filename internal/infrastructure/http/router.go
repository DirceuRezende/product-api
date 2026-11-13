package http

import "net/http"

func NewRouter(productHandler *ProductHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Hello World"))
	})

	mux.HandleFunc("GET /products", productHandler.GetProducts)
	mux.HandleFunc("GET /products/{id}", productHandler.GetProduct)
	mux.HandleFunc("POST /products", productHandler.CreateProduct)
	mux.HandleFunc("PUT /products/{id}", productHandler.AlterProduct)
	mux.HandleFunc("DELETE /products/{id}", productHandler.DeleteProduct)

	return mux
}
