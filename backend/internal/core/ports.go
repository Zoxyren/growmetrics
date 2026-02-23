package core

import (
	"context"

	"esp32/backend/internal/models"
)

type SensorRepository interface {
	SaveSensorData(data models.SensorData) error
	GetAllMetrics() ([]models.SensorData, error)
	GetLatestData(ctx context.Context, deviceID string) (models.SensorData, error)
}

type SensorService interface {
	ProcessIncomingMetric(payload models.SensorPayload, topic string) error
	GetLatestData(ctx context.Context, deviceID string) (models.SensorData, error)
}
