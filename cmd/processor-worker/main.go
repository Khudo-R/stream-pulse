package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Khudo-R/streampulse/internal/broker/rabbitmq"
	"github.com/Khudo-R/streampulse/internal/config"
	"github.com/Khudo-R/streampulse/internal/repository/postgres"
	"github.com/Khudo-R/streampulse/internal/repository/redis"
	"github.com/Khudo-R/streampulse/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()

	dbPool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	log.Println("Worker connected to PostgreSQL")

	repo := postgres.NewEventRepo(dbPool)
	rabbitConn, err := amqp.Dial(cfg.RabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitConn.Close()
	log.Println("Worker connected to RabbitMQ")

	cache, err := redis.NewEventCache(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer cache.Close()
	log.Println("Worker connected to Redis")

	enricher := service.NewStubEnricher()
	consumer, err := rabbitmq.NewConsumer(rabbitConn, repo, enricher, cache)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	if err := consumer.StartConsuming(ctx); err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}

	log.Println("Worker exited gracefully")
}
