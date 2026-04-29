package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Metadata struct {
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
}

type EnrichedData struct {
	Country     string `json:"country"`
	UserSegment string `json:"user_segment"`
}

type Event struct {
	ID           uuid.UUID              `json:"id"`
	UserID       uuid.UUID              `json:"user_id"`
	Type         string                 `json:"type"`
	Payload      map[string]interface{} `json:"payload"`
	Metadata     Metadata               `json:"metadata"`
	EnrichedData *EnrichedData          `json:"enriched_data"`
	CreatedAt    time.Time              `json:"created_at"`
}

func (e *Event) Validate() error {
	if e.UserID == uuid.Nil {
		return fmt.Errorf("user_id is required")
	}
	if e.Type == "" {
		return fmt.Errorf("event type is required")
	}
	return nil
}

type EventRepository interface {
	Save(ctx context.Context, event *Event) error
}

type EventPubslisher interface {
	Publish(ctx context.Context, event *Event) error
}

type Enricher interface {
	GetLocationByIP(ctx context.Context, ip string) (string, error)
	GetUserSegment(ctx context.Context, userID uuid.UUID) (string, error)
}
