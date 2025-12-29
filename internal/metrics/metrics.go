package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	EventsProcessed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "events_processed_total",
			Help: "Total successfully processed events",
		},
	)

	EventsFailed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "events_failed_total",
			Help: "Total failed event processing attempts",
		},
	)

	EventsRetried = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "events_retried_total",
			Help: "Total retried events",
		},
	)

	EventsDLQ = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "events_dlq_total",
			Help: "Total events sent to DLQ",
		},
	)

	ProcessingLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "event_processing_seconds",
			Help:    "Event processing latency",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func Register() {
	prometheus.MustRegister(
		EventsProcessed,
		EventsFailed,
		EventsRetried,
		EventsDLQ,
		ProcessingLatency,
	)
}
