package postgres

import (
	"context"
	"encoding/json"

	"github.com/Khudo-R/streampulse/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepo struct {
	db *pgxpool.Pool
}

func NewEventRepo(db *pgxpool.Pool) *EventRepo {
	return &EventRepo{db: db}
}

func (r *EventRepo) Save(ctx context.Context, e *domain.Event) error {
	query := `INSERT INTO events (id, user_id, type, payload, metadata, enriched_data, created_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7)`

	payload := e.Payload
	if payload == nil {
		payload = make(map[string]interface{})
	}

	metaBytes, _ := json.Marshal(e.Metadata)
	enrichedBytes, _ := json.Marshal(e.EnrichedData)

	_, err := r.db.Exec(ctx, query,
		e.ID,
		e.UserID,
		e.Type,
		payload,
		metaBytes,
		enrichedBytes,
		e.CreatedAt,
	)
	return err
}
