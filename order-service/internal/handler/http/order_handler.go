// order-service/internal/handler/http/order_handler.go
package http

import (
	"encoding/json"
	"net/http"
	"strings" 

	"github.com/Angellieta/order-service/internal/service"
)

type CreateOrderRequest struct {
	ProductID string  `json:"productId"`
	Price     float64 `json:"price"`
	Qty       int     `json:"qty"`
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
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	order, err := h.service.CreateOrder(req.ProductID, req.Price, req.Qty)
	if err != nil {
		// Cek apakah errornya adalah ErrProductNotFound
		if err == service.ErrProductNotFound {
			http.Error(w, err.Error(), http.StatusBadRequest) // Kirim 400 Bad Request
			return
		}
		// Untuk error lainnya, kirim 500 Internal Server Error
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