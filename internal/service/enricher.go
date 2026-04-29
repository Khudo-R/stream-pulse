package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type StubEnricher struct {
	client *http.Client
}

func NewStubEnricher() *StubEnricher {
	return &StubEnricher{
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

type geoAPIResponse struct {
	Status  string `json:"status"`
	Country string `json:"country"`
}

func (e *StubEnricher) GetLocationByIP(ctx context.Context, ip string) (string, error) {
	if ip == "127.0.0.1" || ip == "localhost" || ip == "::1" {
		return "Local", nil
	}
	if ip == "" {
		return "Unknown", nil
	}
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "Error", err
	}
	resp, err := e.client.Do(req)
	if err != nil {
		return "Error", err
	}
	defer resp.Body.Close()

	var geo geoAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&geo); err != nil {
		return "Error", err
	}

	if geo.Status != "success" || geo.Country == "" {
		return "Unknown", nil
	}

	return geo.Country, nil
}

func (e *StubEnricher) GetUserSegment(ctx context.Context, userID uuid.UUID) (string, error) {
	char := userID.String()[0]
	if char >= '0' && char <= '7' {
		return "Newbie", nil
	}
	return "VIP", nil
}
