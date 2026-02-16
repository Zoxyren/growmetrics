package domain

import (
	"fmt"

	"esp32/backend/internal/core"
	"esp32/backend/internal/models"
)

type sensorService struct {
	repo core.SensorRepository
}

func NewSensorService(r core.SensorRepository) core.SensorService {
	return &sensorService{
		repo: r,
	}
}

func (s *sensorService) ProcessIncomingMetric(payload models.SensorPayload, topic string) error {
	customerID, err := s.repo.GetCustomerIDByDeviceID(payload.DeviceID)
	if err != nil {
		return fmt.Errorf("kunde nicht gefunden: %w", err)
	}

	sensorData := models.SensorData{
		DeviceID:    payload.DeviceID,
		CustomerID:  customerID,
		Temperature: payload.Temperature,
		Humidity:    payload.Humidity,
		Pressure:    payload.Pressure,
		Topic:       topic,
	}

	return s.repo.SaveSensorData(sensorData)
}
