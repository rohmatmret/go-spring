// internal/observability/health.go
package observability

type StorageManager struct {
	remote *RemoteStorageConfig
}

func (r *RemoteStorageConfig) CheckHealth() HealthStatus {
	// Add your health check logic here
	// For example, ping the remote storage
	if r.URL == "" {
		return HealthStatusUnhealthy
	}
	return HealthStatusHealthy
}

type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
)

type StorageHealth struct {
	RemoteStorage HealthStatus
}

func (sm *StorageManager) CheckHealth() *StorageHealth {
	return &StorageHealth{
		RemoteStorage: sm.remote.CheckHealth(),
	}
}
