package orders

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

func (h *handler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	order, err := h.service.GetOrderByID(r.Context(), id)
	if err != nil {
		log.Println(err)
		if errors.Is(err, ErrOrderNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, order)
}

func (h *handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {

	// 1. Read the JSON request body
	var tempOrder createOrderParams

	if err := json.Read(r, &tempOrder); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// 2. Call the service to place the order
	orderCreated, err := h.service.PlaceOrder(r.Context(), tempOrder)

	if err != nil {
		log.Println(err)

		if err == ErrProductNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// 3. Return the JSON in an HTTP response
	json.Write(w, http.StatusOK, orderCreated)
}
