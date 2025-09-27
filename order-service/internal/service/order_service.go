// order-service/internal/service/order_service.go

package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Angellieta/order-service/internal/model"
	"github.com/Angellieta/order-service/internal/repository"
	"github.com/Angellieta/order-service/pkg/rabbitmq"
	"github.com/redis/go-redis/v9"
)

var ErrProductNotFound = errors.New("product not found")

// OrderService mendefinisikan interface untuk service order.
type OrderService interface {
	CreateOrder(correlationID, productID string, price float64, qty int) (*model.Order, error)
	GetOrdersByProductID(productID string) ([]model.Order, error)
}

// orderService adalah implementasi dari OrderService.
type orderService struct {
	repo         repository.OrderRepository
	publisher    rabbitmq.EventPublisher
	redisClient  *redis.Client
	ctx          context.Context
}

// NewOrderService membuat instance baru dari orderService dengan semua dependensinya.
func NewOrderService(repo repository.OrderRepository, publisher rabbitmq.EventPublisher, redisClient *redis.Client) OrderService {
	return &orderService{
		repo:         repo,
		publisher:    publisher,
		redisClient:  redisClient,
		ctx:          context.Background(),
	}
}

// CreateOrder berisi logika bisnis untuk membuat pesanan baru.
func (s *orderService) CreateOrder(correlationID, productID string, price float64, qty int) (*model.Order, error) {
	log.Printf("[CorrelationID: %s] Starting to create order for product %s", correlationID, productID)

	// Validasi ke product-service dengan menyertakan Correlation ID
	productURL := fmt.Sprintf("http://product-service:3000/products/%s", productID)
	req, _ := http.NewRequest("GET", productURL, nil)
	req.Header.Set("x-correlation-id", correlationID) // <-- Set header
	
	resp, err := http.DefaultClient.Do(req) // <-- Kirim request
	if err != nil {
		log.Printf("[CorrelationID: %s] Failed to call product-service: %v", correlationID, err)
		return nil, errors.New("internal server error")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrProductNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("product service returned an error")
	}

	totalPrice := price * float64(qty)
	status := "PENDING"
	
	// logika membuat order object
	order := model.Order{
		ProductID:     productID,
		Qty:           qty,
		TotalPrice:    totalPrice,
		Status:        status,
		CorrelationID: correlationID,
		CreatedAt:     time.Now(),
	}

	err = s.repo.CreateOrder(&order)
	if err != nil {
		log.Printf("[CorrelationID: %s] Failed to save order to database: %v", correlationID, err)
		return nil, err
	}
	log.Printf("[CorrelationID: %s] Order %s created successfully", correlationID, order.ID)

	// Menerbitkan event
	go func() {
		err := s.publisher.Publish(order, "order.created")
		if err != nil {
			log.Printf("[CorrelationID: %s] Failed to publish event for order %s: %v", correlationID, order.ID, err)
		}
	}()

	return &order, nil
}

// GetOrdersByProductID mengimplementasikan logika cache-aside untuk mengambil data pesanan.
func (s *orderService) GetOrdersByProductID(productID string) ([]model.Order, error) {
	// Mencoba mengambil data dari Cache (Redis) terlebih dahulu.
	cacheKey := fmt.Sprintf("orders:product:%s", productID)
	cachedOrders, err := s.redisClient.Get(s.ctx, cacheKey).Result()

	if err == nil {
		// Jika ada di cache (Cache Hit)
		log.Println("CACHE HIT: Fetching orders from Redis")
		var orders []model.Order
		err = json.Unmarshal([]byte(cachedOrders), &orders)
		if err != nil {
			return nil, err
		}
		return orders, nil
	}
	
	if err != redis.Nil {
		// Jika terjadi error selain "key not found", log error tapi tetap lanjutkan ke DB.
		log.Printf("Redis error: %v. Fetching from DB as fallback.", err)
	}

	// Jika tidak ada di cache (Cache Miss) atau Redis error
	log.Println("CACHE MISS: Fetching orders from database")
	orders, err := s.repo.GetOrdersByProductID(productID)
	if err != nil {
		return nil, err
	}

	// Simpan hasil dari database ke dalam cache untuk request selanjutnya.
	jsonData, err := json.Marshal(orders)
	if err != nil {
		return nil, err
	}
	// Set cache dengan masa berlaku (TTL) 5 menit
	err = s.redisClient.Set(s.ctx, cacheKey, jsonData, 5*time.Minute).Err()
	if err != nil {
		// Jika gagal menyimpan ke cache, jangan gagalkan request. Cukup log error.
		log.Printf("Failed to set cache for key %s: %v", cacheKey, err)
	}

	return orders, nil
}