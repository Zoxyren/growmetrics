package core

import "esp32/backend/internal/models"

type SensorRepository interface {
	SaveSensorData(data models.SensorData) error
	GetAllMetrics() ([]models.SensorData, error)
	GetCustomerIDByDeviceID(deviceID string) (int, error)
}

type SensorService interface {
	ProcessIncomingMetric(payload models.SensorPayload, topic string) error
}
