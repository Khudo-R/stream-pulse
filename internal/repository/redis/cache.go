package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type EventCache struct {
	client *redis.Client
}

func NewEventCache(redisURL string) (*EventCache, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("invalid redis URL: %w", err)
	}

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return &EventCache{client: client}, nil
}

func (c *EventCache) MarkAsProcessed(ctx context.Context, eventID string) (bool, error) {
	key := fmt.Sprintf("processed_event:%s", eventID)

	success, err := c.client.SetNX(ctx, key, "1", 24*time.Hour).Result()
	if err != nil {
		return false, err
	}

	return !success, nil
}

func (c *EventCache) RemoveMark(ctx context.Context, eventID string) error {
	key := fmt.Sprintf("processed_event:%s", eventID)
	return c.client.Del(ctx, key).Err()
}

func (c *EventCache) Close() error {
	return c.client.Close()
}
