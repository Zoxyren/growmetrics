package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"esp32/backend/internal/models"

	"github.com/redis/go-redis/v9"
)

type RedisAdapter struct {
	rdb *redis.Client
}

func NewRedisAdapter(addr string) *RedisAdapter {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisAdapter{
		rdb: client,
	}
}

func (ra *RedisAdapter) SaveLatestMetrics(ctx context.Context, deviceID string, data models.SensorData) error {
	payload, err := json.Marshal(&data)
	if err != nil {
		slog.Error("Failed to marshal", "error:", err)
		return err
	}
	key := fmt.Sprintf("device:metric:%s", deviceID)
	if err != nil {
		return err
	}
	return ra.rdb.Set(ctx, key, payload, 15*time.Minute).Err()
}

func (ra *RedisAdapter) GetLatestMetrics(ctx context.Context, deviceID string) (models.SensorData, error) {
	key := fmt.Sprintf("device:metric:%s", deviceID)
	data, err := ra.rdb.Get(ctx, key).Result()
	if err != nil {
		return models.SensorData{}, err
	}
	var body models.SensorData
	err = json.Unmarshal([]byte(data), &body)
	if err != nil {
		return models.SensorData{}, err
	}
	return body, nil
}
