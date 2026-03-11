package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Matchmaking metrics collected by Prometheus
var (
	// ActiveQueueSize tracks the number of teams currently in matchmaking queue
	ActiveQueueSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "golobby",
		Subsystem: "matchmaking",
		Name:      "queue_size",
		Help:      "Number of teams currently waiting in the matchmaking queue",
	}, []string{"category"}) // label: POKE | WARKOP

	// MatchFoundTotal counts the total number of successful matches
	MatchFoundTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "golobby",
		Subsystem: "matchmaking",
		Name:      "match_found_total",
		Help:      "Total number of matches successfully created",
	}, []string{"category"})

	// MatchSearchLatency tracks how long it takes to find a match (seconds)
	MatchSearchLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "golobby",
		Subsystem: "matchmaking",
		Name:      "search_latency_seconds",
		Help:      "Duration in seconds from enqueue to match found",
		Buckets:   []float64{1, 2, 5, 10, 15, 20, 30, 45, 60},
	}, []string{"category"})

	// MatchFailedTotal counts times no opponent was found within the timeout
	MatchFailedTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "golobby",
		Subsystem: "matchmaking",
		Name:      "match_failed_total",
		Help:      "Total number of matchmaking attempts that failed to find an opponent",
	}, []string{"category", "reason"})

	// ActiveWebSocketConnections tracks live WebSocket connections
	ActiveWebSocketConnections = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "golobby",
		Name:      "websocket_connections_active",
		Help:      "Number of active WebSocket connections",
	})

	// OCRVerificationTotal counts OCR verification outcomes
	OCRVerificationTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "golobby",
		Subsystem: "ocr",
		Name:      "verification_total",
		Help:      "Total OCR verification outcomes",
	}, []string{"result"}) // label: verified | disputed | error

	// OCRProcessingLatency tracks OCR processing duration
	OCRProcessingLatency = promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "golobby",
		Subsystem: "ocr",
		Name:      "processing_latency_seconds",
		Help:      "Duration in seconds for OCR microservice to return a result",
		Buckets:   []float64{0.5, 1, 2, 5, 10, 20, 30},
	})

	// ReputationPenaltiesTotal tracks how many reputation penalties were applied
	ReputationPenaltiesTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "golobby",
		Subsystem: "reputation",
		Name:      "penalties_total",
		Help:      "Total reputation penalties applied",
	}, []string{"type"}) // label: fraud | ghosting

	// HTTPRequestDuration tracks HTTP handler latency
	HTTPRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "golobby",
		Name:      "http_request_duration_seconds",
		Help:      "Duration of HTTP requests",
		Buckets:   prometheus.DefBuckets,
	}, []string{"method", "route", "status_code"})
)
