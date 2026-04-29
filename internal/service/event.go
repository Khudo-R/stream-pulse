package service

import (
	"context"
	"time"

	"github.com/Khudo-R/streampulse/internal/domain"
	"github.com/google/uuid"
)

type EventService struct {
	publisher domain.EventPubslisher
}

func NewEventService(publisher domain.EventPubslisher) *EventService {
	return &EventService{publisher: publisher}
}

func (s *EventService) CreateEvent(ctx context.Context, event *domain.Event) error {
	if err := event.Validate(); err != nil {
		return err
	}

	if event.ID == uuid.Nil {
		event.ID = uuid.New()
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now()
	}

	return s.publisher.Publish(ctx, event)
}
