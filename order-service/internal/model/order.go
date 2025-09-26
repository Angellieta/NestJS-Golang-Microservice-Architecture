// order-service/internal/model/order.go

package model

import "time"

// Order mendefinisikan struktur tabel 'orders' di database.
type Order struct {
	ID         string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ProductID  string    `json:"productId" gorm:"type:varchar(36);not null"`
	Qty        int       `json:"qty" gorm:"not null"`
	TotalPrice float64   `json:"totalPrice" gorm:"not null"`
	Status     string    `json:"status" gorm:"not null"`
	CreatedAt  time.Time `json:"createdAt"`
}