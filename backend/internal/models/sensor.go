package models

import "time"

type SensorPayload struct {
	ID          int       `json:"id"`
	DeviceID    string    `json:"device_id"`
	Temperature float64   `json:"temperature"`
	Humidity    float64   `json:"humidity"`
	Pressure    float64   `json:"pressure"`
	Topic       string    `json:"topic"`
	CreatedAt   time.Time `json:"created_at"`
}
