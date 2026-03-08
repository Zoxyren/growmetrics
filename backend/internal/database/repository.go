package database

import (
	"context"
	"database/sql"
	"fmt"

	"esp32/backend/internal/models"
)

type Range struct {
	Min float64
	Max float64
}

var MetricRanges = map[string]Range{
	"temperature": {Min: 0.1, Max: 100.0},
	"humidity":    {Min: 0.0, Max: 100.0},
	"pressure":    {Min: 300.0, Max: 1100.0},
}

type DatabaseAdapter struct {
	db *sql.DB
}

func NewDatabaseAdapter(db *sql.DB) *DatabaseAdapter {
	return &DatabaseAdapter{db: db}
}

func (da *DatabaseAdapter) GetAllMetrics() ([]models.SensorData, error) {
	rows, err := da.db.Query(`
        SELECT id, device_id, temperature, humidity, pressure, topic, created_at 
        FROM sensor_data 
        ORDER BY created_at DESC 
        LIMIT 100`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []models.SensorData
	for rows.Next() {
		var v models.SensorData
		err := rows.Scan(&v.ID, &v.DeviceID, &v.Temperature, &v.Humidity, &v.Pressure, &v.Topic, &v.CreatedAt)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, v)
	}
	return metrics, nil
}

func (da *DatabaseAdapter) GetCustomerIDByDeviceID(deviceID string) (int, error) {
	var customerID int
	err := da.db.QueryRow("SELECT customer_id FROM devices WHERE id = $1", deviceID).Scan(&customerID)
	return customerID, err
}

func (da *DatabaseAdapter) SaveSensorData(data models.SensorData) error {
	metricstoCheck := map[string]float64{
		"temperature": data.Temperature,
		"humidity":    data.Humidity,
		"pressure":    data.Pressure,
	}
	for name, value := range metricstoCheck {
		if r, ok := MetricRanges[name]; ok {
			if value < r.Min || value > r.Max {
				return fmt.Errorf("validierung fehlgeschlagen: %s (%.2f) außerhalb Bereich", name, value)
			}
		}
	}

	query := `
        INSERT INTO sensor_data (device_id, temperature, humidity, pressure, topic)
        VALUES ($1, $2, $3, $4, $5)`

	_, err := da.db.Exec(query,
		data.DeviceID,
		//	data.CustomerID,
		data.Temperature,
		data.Humidity,
		data.Pressure,
		data.Topic,
	)
	return err
}

func (da *DatabaseAdapter) GetLatestData(ctx context.Context, deviceID string) (models.SensorData, error) {
	var data models.SensorData
	err := da.db.QueryRow("SELECT device_id, temperature, humidity, pressure, topic FROM SensorData WHERE device_id = $1 ORDER BY created_at DESC LIMIT 1;", deviceID).Scan(
		&data.DeviceID,
		&data.Temperature,
		&data.Humidity,
		&data.Pressure,
		&data.Topic,
	)
	return data, err
}
