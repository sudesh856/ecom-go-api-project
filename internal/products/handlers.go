package products

import (
	"log"
	"net/http"
	"strconv"

	repo "github.com/sudesh856/ecom-go-api-project/internal/adaptors/postgresql/sqlc"
	"github.com/sudesh856/ecom-go-api-project/internal/json"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{

		service: service,
	}
}

func (h *handler) ListProducts(w http.ResponseWriter, r *http.Request) {

	products, err := h.service.ListProducts(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	
	json.Write(w, http.StatusOK, products)

}

func (h *handler) FindProduct(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	productID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}

	another_product, err := h.service.FindProduct(r.Context(), productID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.Write(w, http.StatusOK, another_product)
}

func (h *handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var arg repo.CreateProductParams

	if err := json.Read(r, &arg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	product, err := h.service.CreateProduct(r.Context(), arg)
    if err != nil {
        log.Println(err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.Write(w, http.StatusCreated, product)
}

