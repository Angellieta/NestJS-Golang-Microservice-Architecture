// order-service/pkg/rabbitmq/publisher.go
package rabbitmq

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// EventPublisher adalah interface untuk mengirim event.
type EventPublisher interface {
	Publish(body interface{}, routingKey string) error
}

type amqpPublisher struct {
	conn *amqp.Connection
}

// NewPublisher membuat koneksi ke RabbitMQ dan sebuah exchange.
func NewPublisher() (EventPublisher, error) {
	amqpURL := os.Getenv("RABBITMQ_URL")
	if amqpURL == "" {
		log.Fatal("RABBITMQ_URL environment variable not set")
	}

	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	// Buat channel
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	// Deklarasikan exchange tempat yang akan mengirim pesan
	// Menggunakan 'topic' exchange untuk fleksibilitas routing
	err = ch.ExchangeDeclare(
		"orders_exchange", // nama exchange
		"topic",           // tipe
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		return nil, err
	}
	log.Println("RabbitMQ exchange 'orders_exchange' declared")

	return &amqpPublisher{conn: conn}, nil
}

// Publish mengirim pesan ke RabbitMQ.
func (p *amqpPublisher) Publish(body interface{}, routingKey string) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Ubah body (struct) menjadi JSON
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("Publishing message to exchange 'orders_exchange' with routing key '%s'", routingKey)

	return ch.PublishWithContext(ctx,
		"orders_exchange", // exchange
		routingKey,        // routing key (e.g., "order.created")
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonBody,
		},
	)
}