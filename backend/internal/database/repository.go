package database

import (
	"database/sql"

	"esp32/backend/internal/models"
)

type DatabaseAdapter struct {
	db *sql.DB
}

func NewDatabaseAdapter(db *sql.DB) *DatabaseAdapter {
	return &DatabaseAdapter{db: db}
}

func (da *DatabaseAdapter) GetAllMetrics() ([]models.SensorData, error) {
	rows, err := da.db.Query(`
        SELECT id, device_id, customer_id, temperature, humidity, pressure, topic, created_at 
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
		err := rows.Scan(&v.ID, &v.DeviceID, &v.CustomerID, &v.Temperature, &v.Humidity, &v.Pressure, &v.Topic, &v.CreatedAt)
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
	query := `
        INSERT INTO sensor_data (device_id, customer_id, temperature, humidity, pressure, topic)
        VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := da.db.Exec(query,
		data.DeviceID,
		data.CustomerID,
		data.Temperature,
		data.Humidity,
		data.Pressure,
		data.Topic,
	)
	return err
}
