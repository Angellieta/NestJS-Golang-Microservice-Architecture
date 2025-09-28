// order-service/internal/service/order_service_test.go
package service

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/Angellieta/order-service/internal/model"
	"github.com/redis/go-redis/v9"
)

// --- Mock Dependencies ---

type mockOrderRepository struct {
	err error
}
func (m *mockOrderRepository) CreateOrder(order *model.Order) error { return m.err }
func (m *mockOrderRepository) GetOrdersByProductID(productID string) ([]model.Order, error) { return nil, nil }

type mockEventPublisher struct {
	err error
}
func (m *mockEventPublisher) Publish(body interface{}, routingKey string) error { return m.err }

type mockHttpClient struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHttpClient) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}
// --- Test Function ---

func TestCreateOrder(t *testing.T) {
	// Setup mock HTTP Client
	originalClient := http.DefaultClient
	defer func() { http.DefaultClient = originalClient }() // Kembalikan setelah tes selesai
	mockClient := &mockHttpClient{}
	http.DefaultClient = &http.Client{
		Transport: mockClient,
	}

	// Skenario 1: Sukses
	t.Run("should create order successfully", func(t *testing.T) {
		// Arrange
		mockClient.RoundTripFunc = func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
			}, nil
		}
		mockRepo := &mockOrderRepository{err: nil}
		mockPub := &mockEventPublisher{err: nil}
		dummyRedis := redis.NewClient(&redis.Options{})
		orderService := NewOrderService(mockRepo, mockPub, dummyRedis)

		// Act
		order, err := orderService.CreateOrder("test-correlation-id", "a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d", 100, 2)

		// Assert
		if err != nil { t.Errorf("expected no error, but got %v", err) }
		if order == nil { t.Errorf("expected order not to be nil") }
	})

	// Skenario 2: Gagal
	t.Run("should return error when repository fails", func(t *testing.T) {
		// Arrange
		mockClient.RoundTripFunc = func(req *http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewBufferString(`{}`))} , nil
		}
		
		expectedErr := errors.New("database error")
		mockRepo := &mockOrderRepository{err: expectedErr}
		mockPub := &mockEventPublisher{err: nil}
		dummyRedis := redis.NewClient(&redis.Options{})
		orderService := NewOrderService(mockRepo, mockPub, dummyRedis)

		// Act
		_, err := orderService.CreateOrder("test-correlation-id", "a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d", 100, 2)

		// Assert
		if err == nil { t.Errorf("expected an error, but got nil") }
	})
}