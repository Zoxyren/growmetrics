package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
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

func (ra *RedisAdapter) GetAllDevices(ctx context.Context) ([]string, error) {
	var allKeys []string
	cursor := uint64(0)
	prefix := "active:device:*"

	for {
		keys, nextcursor, err := ra.rdb.Scan(ctx, cursor, prefix, 100).Result()
		if err != nil {
			return nil, fmt.Errorf("fehler beim scannen der redis keys: %w", err)
		}

		allKeys = append(allKeys, keys...)

		cursor = nextcursor

		if cursor == 0 {
			break
		}
	}

	return ra.extractDeviceIDs(allKeys), nil
}

func (ra *RedisAdapter) extractDeviceIDs(keys []string) []string {
	ids := make([]string, 0, len(keys))
	prefix := "active:device:"

	for _, key := range keys {
		// Wir schneiden den Teil "active:device:" einfach ab
		id := strings.TrimPrefix(key, prefix)
		ids = append(ids, id)
	}
	return ids
}

func (ra *RedisAdapter) SetDeviceActive(ctx context.Context, deviceID string) error {
	key := fmt.Sprintf("active:device:%s", deviceID)
	// Wir setzen den Wert "1" mit 10 Min Ablaufzeit
	return ra.rdb.Set(ctx, key, "1", 10*time.Minute).Err()
}
