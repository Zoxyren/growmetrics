package database

import (
	"database/sql"
	"log/slog"

	"esp32/backend/internal/models"
)

type DatabaseAdapter struct {
	db *sql.DB
}

func NewDatabaseAdapter(db *sql.DB) *DatabaseAdapter {
	return &DatabaseAdapter{
		db: db,
	}
}

func (da *DatabaseAdapter) GetAllMetrics() ([]models.SensorPayload, error) {
	metrics := []models.SensorPayload{}
	// Anpassungen an dem SQL Statements für den Response
	query, err := da.db.Query(`SELECT id, device_id, temperature, humidity, pressure, topic, created_at FROM sensor_data`)
	if err != nil {
		return nil, err
	}
	defer query.Close()
	for query.Next() {
		var v models.SensorPayload
		err = query.Scan(&v.ID, &v.DeviceID, &v.Temperature, &v.Humidity, &v.Pressure, &v.Topic, &v.CreatedAt)
		if err != nil {
			slog.Error("Failed to recieve metrics", "error", err)
		}
		if err := query.Err(); err != nil {
			return nil, err
		}
		metrics = append(metrics, v)
	}

	return metrics, nil
}
