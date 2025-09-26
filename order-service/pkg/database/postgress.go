// order-service/pkg/database/postgres.go

package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewConnection membuat dan mengembalikan koneksi database GORM.
func NewConnection() (*gorm.DB, error) {
	// Mengambil URL database dari environment variable yang sudah set di docker-compose.yml
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	// Membuka koneksi ke database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Database connection successfully established")
	return db, nil
}