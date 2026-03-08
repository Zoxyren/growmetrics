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
	redis core.StatusRepository
}

func NewSensorService(r core.SensorRepository, c core.CacheRepository, red core.StatusRepository) core.SensorService {
	return &sensorService{
		repo:  r,
		cache: c,
		redis: red,
	}
}

func (s *sensorService) ProcessIncomingMetric(payload models.SensorPayload, topic string) error {
	sensorData := models.SensorData{
		DeviceID:    payload.DeviceID,
		Temperature: payload.Temperature,
		Humidity:    payload.Humidity,
		Pressure:    payload.Pressure,
		Topic:       topic,
	}
	err := s.cache.SaveLatestMetrics(context.Background(), sensorData.DeviceID, sensorData)
	if err != nil {
		slog.Info("failed to save data in cache", "error", err)
	}
	_ = s.redis.SetDeviceActive(context.Background(), payload.DeviceID)

	return s.repo.SaveSensorData(sensorData)
}

func (s *sensorService) GetLatestData(ctx context.Context, deviceID string) (models.SensorData, error) {
	data, err := s.cache.GetLatestMetrics(ctx, deviceID)
	if err == nil {
		return data, nil
	}
	dbData, err := s.repo.GetLatestData(ctx, deviceID)
	if err != nil {
		return models.SensorData{}, fmt.Errorf("sensor data not found: %w", err)
	}
	_ = s.cache.SaveLatestMetrics(ctx, deviceID, dbData)

	return dbData, nil
}

func (s *sensorService) GetActiveDevices(ctx context.Context) ([]string, error) {
	return s.redis.GetAllDevices(ctx)
}
