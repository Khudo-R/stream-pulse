package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Khudo-R/streampulse/internal/domain"
	"github.com/Khudo-R/streampulse/internal/repository/redis"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn     *amqp.Connection
	ch       *amqp.Channel
	repo     domain.EventRepository
	enricher domain.Enricher
	cache    *redis.EventCache
}

func NewConsumer(conn *amqp.Connection, repo domain.EventRepository, enricher domain.Enricher, cache *redis.EventCache) (*Consumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &Consumer{
		conn:     conn,
		ch:       ch,
		repo:     repo,
		enricher: enricher,
		cache:    cache,
	}, nil
}

func (c *Consumer) StartConsuming(ctx context.Context) error {
	msgs, err := c.ch.Consume(
		"events_queue",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	log.Println("👷 Worker started. Waiting for messages in 'events_queue'...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Worker stopping due to context cancellation...")
			return nil

		case d, ok := <-msgs:
			if !ok {
				return fmt.Errorf("channel closed")
			}

			var event domain.Event
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("Failed to parse message: %v. Raw body: %s", err, string(d.Body))
				d.Ack(false)
				continue
			}

			isDuplicate, err := c.cache.MarkAsProcessed(ctx, event.ID.String())
			if err != nil {
				log.Printf("⚠️ Redis error: %v", err)
				d.Nack(false, true)
				continue
			}
			if isDuplicate {
				log.Printf("Duplicate event detected and skipped: %s", event.ID)
				d.Ack(false)
				continue
			}

			if event.EnrichedData == nil {
				event.EnrichedData = &domain.EnrichedData{}
			}

			country, _ := c.enricher.GetLocationByIP(ctx, event.Metadata.IP)
			segment, _ := c.enricher.GetUserSegment(ctx, event.UserID)

			event.EnrichedData.Country = country
			event.EnrichedData.UserSegment = segment

			if err := c.repo.Save(ctx, &event); err != nil {
				log.Printf("Failed to save event %s to DB: %v", event.ID, err)
				d.Nack(false, true)
				continue
			}

			log.Printf("✅ Enriched & saved event: %s (Country: %s, Segment: %s)", event.ID, country, segment)
			d.Ack(false)
		}
	}
}

func (c *Consumer) Close() {
	c.ch.Close()
}
