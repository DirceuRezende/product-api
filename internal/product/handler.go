package product

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Request é o corpo JSON de criação/alteração.
type Request struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       uint64 `json:"price"`
	Stock       uint64 `json:"stock"`
}

// Response é a representação JSON de um produto.
type Response struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       uint64 `json:"price"`
	Stock       uint64 `json:"stock"`
}

func newResponse(p *Product) Response {
	return Response{ID: p.ID, Name: p.Name, Description: p.Description, Price: p.Price, Stock: p.Stock}
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
	switch {
	case errors.Is(err, ErrNotFound):
		writeError(w, http.StatusNotFound, ErrNotFound.Error())
	case errors.Is(err, ErrInvalidID):
		writeError(w, http.StatusBadRequest, ErrInvalidID.Error())
	default:
		log.Printf("internal error: %v", err)
		writeError(w, http.StatusInternalServerError, "erro interno")
	}
}

func parseID(r *http.Request) (uint64, bool) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id == 0 {
		return 0, false
	}
	return id, true
}

func decodeRequest(w http.ResponseWriter, r *http.Request) (Input, bool) {
	var req Request

	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return Input{}, false
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name é obrigatório")
		return Input{}, false
	}
	return Input{Name: req.Name, Description: req.Description, Price: req.Price, Stock: req.Stock}, true
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, ErrInvalidID.Error())
		return
	}

	p, err := h.service.Get(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newResponse(p))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.List(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	res := make([]Response, 0, len(products))
	for _, p := range products {
		res = append(res, newResponse(p))
	}

	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	in, ok := decodeRequest(w, r)
	if !ok {
		return
	}

	p, err := h.service.Create(r.Context(), in)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Location", "/products/"+strconv.FormatUint(p.ID, 10))
	writeJSON(w, http.StatusCreated, newResponse(p))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, ErrInvalidID.Error())
		return
	}

	in, ok := decodeRequest(w, r)
	if !ok {
		return
	}

	p, err := h.service.Update(r.Context(), id, in)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newResponse(p))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, ErrInvalidID.Error())
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
