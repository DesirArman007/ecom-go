package products

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/DesirArman007/ecom/internal/json"
	"github.com/go-chi/chi/v5"
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

	// 1. Call the service to list the products

	products, err := h.service.ListProducts(r.Context())

	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// 2. Return JSON in an HTTP response
	json.Write(w, http.StatusOK, products)
}

func (h *handler) FindProductByID(w http.ResponseWriter, r *http.Request) {

	// 1. Get the product id from request url
	idString := chi.URLParam(r, "id")

	productID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// 2. Call the service to find the product if exits by the id
	product, err := h.service.FindProductByID(r.Context(), productID)

	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// 3. Return JSON in an HTTP response
	json.Write(w, http.StatusOK, product)
}

func (h *handler) CreateProduct(w http.ResponseWriter, r *http.Request) {

	// 1. Read the JSON request body
	var tempProduct createProductParams

	if err := json.Read(r, &tempProduct); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// 2. Call the service to create the product
	productCreated, err := h.service.CreateProduct(r.Context(), tempProduct)

	if err != nil {
		log.Println(err)

		if errors.Is(err, ErrProductNameRequired) ||
			errors.Is(err, ErrPriceInvalid) ||
			errors.Is(err, ErrQuantityInvalid) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	//3. Return the JSON in http response
	json.Write(w, http.StatusCreated, productCreated)

}
