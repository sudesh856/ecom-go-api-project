package products

import (
	"log"
	"net/http"
	"strconv"

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
