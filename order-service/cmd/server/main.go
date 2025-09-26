// order-service/cmd/server/main.go
package main

import (
	"log"
	"net/http"
	"strings"

	orderHandler "github.com/Angellieta/order-service/internal/handler/http"
	"github.com/Angellieta/order-service/internal/model"
	"github.com/Angellieta/order-service/internal/repository"
	"github.com/Angellieta/order-service/internal/service"
	"github.com/Angellieta/order-service/pkg/database"
	"github.com/Angellieta/order-service/pkg/rabbitmq"
	"github.com/Angellieta/order-service/pkg/redis"
)

func main() {
	// Setup koneksi
	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	redisClient, err := redis.NewClient()
	if err != nil {
		log.Fatalf("could not connect to redis: %v", err)
	}

	publisher, err := rabbitmq.NewPublisher()
	if err != nil {
		log.Fatalf("could not setup rabbitmq publisher: %v", err)
	}

	// Migrasi DB
	log.Println("Running database migration...")
	db.AutoMigrate(&model.Order{})

	// Dependency Injection
	orderRepo := repository.NewOrderRepository(db)
	orderSvc := service.NewOrderService(orderRepo, publisher, redisClient) // <-- Suntikkan redisClient
	orderHdlr := orderHandler.NewOrderHandler(orderSvc)

	// Mendaftarkan routes/endpoints
	mux := http.NewServeMux()
	mux.HandleFunc("/orders", orderHdlr.CreateOrder)
	mux.HandleFunc("/orders/product/", orderHdlr.GetOrdersByProductID) // <-- Route baru

	// Router sederhana untuk memvalidasi method HTTP
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/orders" {
			if r.Method == http.MethodPost {
				mux.ServeHTTP(w, r)
				return
			}
		} else if strings.HasPrefix(r.URL.Path, "/orders/product/") {
			if r.Method == http.MethodGet {
				mux.ServeHTTP(w, r)
				return
			}
		}
		http.Error(w, "Method not allowed or invalid path", http.StatusMethodNotAllowed)
	})

	// Mulai Server
	port := ":8080"
	log.Printf("Starting server on port %s\n", port)
	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}