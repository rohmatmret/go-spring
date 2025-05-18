package observability

import (
	"github.com/prometheus/client_golang/prometheus"
)

type CacheMetrics struct {
	hits   prometheus.Counter
	misses prometheus.Counter
	errors prometheus.Counter
	size   prometheus.Gauge
}

func NewCacheMetrics() *CacheMetrics {
	return &CacheMetrics{
		hits: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		}),
		misses: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		}),
		errors: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "cache_errors_total",
			Help: "Total number of cache errors",
		}),
		size: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "cache_size",
			Help: "Current size of the cache",
		}),
	}
}

func (m *CacheMetrics) RecordHit() {
	m.hits.Inc()
}

func (m *CacheMetrics) RecordMiss() {
	m.misses.Inc()
}

func (m *CacheMetrics) RecordError() {
	m.errors.Inc()
}

func (m *CacheMetrics) UpdateSize(size int) {
	m.size.Set(float64(size))
}
