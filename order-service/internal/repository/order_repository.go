// order-service/internal/repository/order_repository.go

package repository

import (
	"github.com/Angellieta/order-service/internal/model"
	"gorm.io/gorm"
)

// OrderRepository adalah interface untuk operasi database order.
type OrderRepository interface {
	CreateOrder(order *model.Order) error
	GetOrdersByProductID(productID string) ([]model.Order, error) 
}

// orderRepository adalah implementasi dari OrderRepository.
type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository membuat instance baru dari orderRepository.
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

// CreateOrder menyimpan sebuah order ke dalam database.
func (r *orderRepository) CreateOrder(order *model.Order) error {
	result := r.db.Create(order)
	return result.Error
}

// GetOrdersByProductID mengambil semua order dari database untuk productID tertentu.
func (r *orderRepository) GetOrdersByProductID(productID string) ([]model.Order, error) {
	var orders []model.Order
	result := r.db.Where("product_id = ?", productID).Find(&orders)
	return orders, result.Error
}