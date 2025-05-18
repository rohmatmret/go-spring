package observability

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Example with InfluxDB
type RemoteStorageConfig struct {
	Type      string // "influxdb", "timescaledb", etc.
	URL       string
	Database  string
	Retention string
}

// internal/observability/metrics.go
type MetricsConfig struct {
	// Local storage retention
	LocalRetention time.Duration
	// Remote storage configuration
	RemoteStorage *RemoteStorageConfig
	// Scrape interval
	ScrapeInterval time.Duration
}

func NewMetricsConfig() *MetricsConfig {
	return &MetricsConfig{
		LocalRetention: 15 * 24 * time.Hour, // 15 days
		ScrapeInterval: 15 * time.Second,
		RemoteStorage: &RemoteStorageConfig{
			Type:      "influxdb",
			URL:       "http://influxdb:8086",
			Database:  "metrics",
			Retention: "30d",
		},
	}
}

var (
	// HTTPRequestDuration tracks HTTP request duration
	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	// HTTPRequestsTotal tracks total HTTP requests
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// CacheHits tracks cache hit rate
	CacheHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
		[]string{"cache", "operation"},
	)

	// CacheMisses tracks cache misses
	CacheMisses = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
		[]string{"cache", "operation"},
	)

	// ServiceMethodDuration tracks service method duration
	ServiceMethodDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "service_method_duration_seconds",
			Help:    "Duration of service method calls in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "method"},
	)
)
