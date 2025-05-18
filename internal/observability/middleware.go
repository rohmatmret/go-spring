package observability

import (
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"
)

// MetricsMiddleware adds Prometheus metrics to HTTP requests
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer that captures the status code
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		// Record metrics
		duration := time.Since(start).Seconds()
		HTTPRequestDuration.WithLabelValues(r.Method, r.URL.Path, string(rw.statusCode)).Observe(duration)
		HTTPRequestsTotal.WithLabelValues(r.Method, r.URL.Path, string(rw.statusCode)).Inc()
	})
}

// TracingMiddleware adds OpenTelemetry tracing to HTTP requests
func TracingMiddleware(tracer *Tracer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, span := tracer.StartSpan(r.Context(), r.URL.Path)
			defer span.End()

			// Add request attributes to span
			span.SetAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.url", r.URL.String()),
				attribute.String("http.user_agent", r.UserAgent()),
			)

			// Create a response writer that captures the status code
			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(rw, r.WithContext(ctx))

			// Add response attributes to span
			span.SetAttributes(
				attribute.Int("http.status_code", rw.statusCode),
			)
		})
	}
}

// responseWriter is a wrapper around http.ResponseWriter that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
