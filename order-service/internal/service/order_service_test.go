// order-service/internal/service/order_service_test.go
package service

import (
	"errors"
	"testing"

	"github.com/Angellieta/order-service/internal/model"
)

// --- Mock Dependencies ---

// Mock untuk OrderRepository
type mockOrderRepository struct {
	// Jika CreateOrder dipanggil, akan mengembalikan error ini.
	// Jika nil, artinya sukses.
	err error
}

func (m *mockOrderRepository) CreateOrder(order *model.Order) error {
	return m.err
}

func (m *mockOrderRepository) GetOrdersByProductID(productID string) ([]model.Order, error) {
	return nil, nil
}

// Mock untuk EventPublisher
type mockEventPublisher struct {
	err error
}

func (m *mockEventPublisher) Publish(body interface{}, routingKey string) error {
	return m.err
}

// --- Test Function ---

func TestCreateOrder(t *testing.T) {
	// Skenario 1: Semuanya berjalan sukses
	t.Run("should create order successfully", func(t *testing.T) {
		// Arrange: Menyiapkan semua mock dan service
		mockRepo := &mockOrderRepository{err: nil} // Tidak ada error dari repo
		mockPub := &mockEventPublisher{err: nil}    // Tidak ada error dari publisher
		
		// Buat service dengan dependensi mock
		orderService := NewOrderService(mockRepo, mockPub, nil) // redisClient bisa nil karena tidak dipakai di CreateOrder

		// Act: Jalankan method yang ingin dites
		order, err := orderService.CreateOrder("product-123", 100, 2)

		// Assert: Periksa hasilnya
		if err != nil {
			t.Errorf("expected no error, but got %v", err)
		}
		if order == nil {
			t.Errorf("expected order not to be nil")
		}
		if order.TotalPrice != 200 {
			t.Errorf("expected total price to be 200, but got %f", order.TotalPrice)
		}
	})

	// Skenario 2: Gagal saat menyimpan ke database
	t.Run("should return error when repository fails", func(t *testing.T) {
		// Arrange
		expectedErr := errors.New("database error")
		mockRepo := &mockOrderRepository{err: expectedErr} // Atur repo agar gagal
		mockPub := &mockEventPublisher{err: nil}
		orderService := NewOrderService(mockRepo, mockPub, nil)

		// Act
		order, err := orderService.CreateOrder("product-123", 100, 2)

		// Assert
		if err == nil {
			t.Errorf("expected an error, but got nil")
		}
		if err != expectedErr {
			t.Errorf("expected error '%v', but got '%v'", expectedErr, err)
		}
		if order != nil {
			t.Errorf("expected order to be nil on failure")
		}
	})
}