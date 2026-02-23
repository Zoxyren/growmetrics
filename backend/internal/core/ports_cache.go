package core

import (
	"context"

	"esp32/backend/internal/models"
)

type CacheRepository interface {
	SaveLatestMetrics(ctx context.Context, deviceID string, data models.SensorData) error
	GetLatestMetrics(ctx context.Context, deviceID string) (models.SensorData, error)
}
