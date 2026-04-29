package service

import (
	"context"
	"testing"

	"github.com/Khudo-R/streampulse/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockPublisher struct{}

func (p *mockPublisher) Publish(ctx context.Context, event *domain.Event) error {
	return nil
}

func TestCreateEvent_PreservesID(t *testing.T) {
	pub := &mockPublisher{}
	s := NewEventService(pub)

	fixedID := uuid.New()
	event := &domain.Event{
		ID:     fixedID,
		UserID: uuid.New(),
		Type:   "test",
	}

	err := s.CreateEvent(context.Background(), event)
	assert.NoError(t, err)
	assert.Equal(t, fixedID, event.ID, "Event ID should be preserved if provided")
}
