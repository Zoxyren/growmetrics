package domain

import (
	"context"
	"fmt"
	"log/slog"

	"esp32/backend/internal/core"
	"esp32/backend/internal/models"
)

type sensorService struct {
	repo  core.SensorRepository
	cache core.CacheRepository
}

func NewSensorService(r core.SensorRepository, c core.CacheRepository) core.SensorService {
	return &sensorService{
		repo:  r,
		cache: c,
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
	err = s.cache.SaveLatestMetrics(context.Background(), sensorData.DeviceID, sensorData)
	if err != nil {
		slog.Info("failed to save data in cache", "error", err)
	}

	return s.repo.SaveSensorData(sensorData)
}
