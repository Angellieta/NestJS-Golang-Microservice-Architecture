// order-service/internal/handler/http/order_handler.go
package http

import (
	"encoding/json"
	"net/http"
	"fmt"
	"strings" 

	"github.com/Angellieta/order-service/internal/service"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// DTO untuk request pembuatan order dengan aturan validasi
type CreateOrderRequest struct {
	ProductID string  `json:"productId" validate:"required,uuid"`
	Price     float64 `json:"price"     validate:"required,gt=0"`
	Qty       int     `json:"qty"       validate:"required,gte=1"`
}


type OrderHandler struct {
	service service.OrderService
}

func NewOrderHandler(svc service.OrderService) *OrderHandler {
	return &OrderHandler{
		service: svc,
	}
}

// CreateOrder adalah handler untuk POST /orders juga menangani ErrProductNotFound
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	correlationID := r.Header.Get("x-correlation-id")
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 2. Jalankan validasi pada struct request
	err := validate.Struct(req)
	if err != nil {
		// Jika validasi gagal, kirim error yang detail
		validationErrors := err.(validator.ValidationErrors)
		errorMsg := fmt.Sprintf("Validation error: %s", validationErrors)
		http.Error(w, errorMsg, http.StatusBadRequest)
		return
	}

	// Jika validasi berhasil, lanjutkan ke service
	order, err := h.service.CreateOrder(correlationID, req.ProductID, req.Price, req.Qty)
	if err != nil {
		if err == service.ErrProductNotFound {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// GetOrdersByProductID adalah handler untuk GET /orders/product/{id}
func (h *OrderHandler) GetOrdersByProductID(w http.ResponseWriter, r *http.Request) {
	// Ekstrak productID dari URL, contoh: /orders/product/123 -> 123
	productID := strings.TrimPrefix(r.URL.Path, "/orders/product/")

	// Panggil service untuk mendapatkan data (dari cache atau DB)
	orders, err := h.service.GetOrdersByProductID(productID)
	if err != nil {
		http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}

	// Kirim respons sukses (200 OK) dengan data order dalam format JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}