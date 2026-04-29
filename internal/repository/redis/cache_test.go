package redis

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestEventCache_MarkAndRemove(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	cache := &EventCache{client: client}

	ctx := context.Background()
	eventID := "test-event-123"

	// First mark
	isDuplicate, err := cache.MarkAsProcessed(ctx, eventID)
	assert.NoError(t, err)
	assert.False(t, isDuplicate)

	// Second mark (should be duplicate)
	isDuplicate, err = cache.MarkAsProcessed(ctx, eventID)
	assert.NoError(t, err)
	assert.True(t, isDuplicate)

	// Remove mark
	err = cache.RemoveMark(ctx, eventID)
	assert.NoError(t, err)

	// Third mark (should NOT be duplicate anymore)
	isDuplicate, err = cache.MarkAsProcessed(ctx, eventID)
	assert.NoError(t, err)
	assert.False(t, isDuplicate)
}
