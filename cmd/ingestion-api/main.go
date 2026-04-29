package main

import (
	"log"
	"net/http"

	"github.com/Khudo-R/streampulse/internal/broker/rabbitmq"
	"github.com/Khudo-R/streampulse/internal/config"
	"github.com/Khudo-R/streampulse/internal/service"
	api "github.com/Khudo-R/streampulse/internal/transport/http"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	cfg := config.Load()

	conn, err := amqp.Dial(cfg.RabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	log.Println("Connected to RabbitMQ")

	publisher, err := rabbitmq.NewPublisher(conn)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ publisher: %v", err)
	}
	defer publisher.Close()

	svc := service.NewEventService(publisher)
	handler := api.NewHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/events", handler.PostEvent)

	log.Printf("🚀 Ingestion API started on port %s", cfg.AppPort)
	log.Fatal(http.ListenAndServe(":"+cfg.AppPort, mux))

}
