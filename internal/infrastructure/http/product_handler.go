package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	product "products-api/internal/application/product"
	domain "products-api/internal/domain/product"
)

type ProductHandler struct {
	service *product.ProductService
}

func NewProductHandler(service *product.ProductService) *ProductHandler {
	return &ProductHandler{service}
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if body == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResponse{Error: msg})
}

func handleError(w http.ResponseWriter, err error) {
	if errors.Is(err, domain.ErrNotFound) {
		writeError(w, http.StatusNotFound, "produto não encontrado")
		return
	}
	internalError(w, err)
}

func internalError(w http.ResponseWriter, err error) {
	log.Printf("internal error: %v", err)
	writeError(w, http.StatusInternalServerError, "erro interno")
}

func parseID(r *http.Request) (uint64, bool) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id == 0 {
		return 0, false
	}
	return id, true
}

func decodeRequest(w http.ResponseWriter, r *http.Request) (product.ProductRequest, bool) {
	var req product.ProductRequest

	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return req, false
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name é obrigatório")
		return req, false
	}
	return req, true
}

func (p *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	res, err := p.service.GetProduct(id)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func (p *ProductHandler) GetProducts(w http.ResponseWriter, _ *http.Request) {
	res, err := p.service.GetProducts()
	if err != nil {
		internalError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func (p *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeRequest(w, r)
	if !ok {
		return
	}

	res, err := p.service.CreateProduct(req)
	if err != nil {
		internalError(w, err)
		return
	}

	w.Header().Set("Location", "/products/"+strconv.FormatUint(res.ID, 10))
	writeJSON(w, http.StatusCreated, res)
}

func (p *ProductHandler) AlterProduct(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	req, ok := decodeRequest(w, r)
	if !ok {
		return
	}

	res, err := p.service.AlterProduct(id, req)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func (p *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	if err := p.service.DeleteProduct(id); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
