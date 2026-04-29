package http

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"

	"github.com/Khudo-R/streampulse/internal/domain"
	"github.com/Khudo-R/streampulse/internal/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var eventsProcessed = promauto.NewCounter(prometheus.CounterOpts{
	Name: "streampulse_incoming_events_total",
	Help: "The total number of incoming events",
})

type Handler struct {
	svc *service.EventService
}

func NewHandler(svc *service.EventService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) PostEvent(w http.ResponseWriter, r *http.Request) {
	eventsProcessed.Inc()
	var event domain.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if event.Metadata.IP == "" {
		forwarded := r.Header.Get("X-Forwarded-For")
		if forwarded != "" {
			event.Metadata.IP = strings.Split(forwarded, ",")[0]
		} else {
			host, _, err := net.SplitHostPort(r.RemoteAddr)
			if err == nil {
				event.Metadata.IP = host
			} else {
				event.Metadata.IP = r.RemoteAddr
			}
		}
	}

	if err := h.svc.CreateEvent(r.Context(), &event); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
