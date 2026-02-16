package models

import "time"

type Customer struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	CreatedAt time.Time `json:"created_at"`
}

type Device struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	CustomerID int       `json:"customer_id" gorm:"index;not null"`
	DeviceName string    `json:"device_name"`
	CreatedAt  time.Time `json:"created_at"`

	Customer Customer `json:"-" gorm:"foreignKey:CustomerID"`
}

type SensorData struct {
	ID          int64     `json:"id" gorm:"primaryKey"`
	DeviceID    string    `json:"device_id" gorm:"index;not null"`
	CustomerID  int       `json:"customer_id" gorm:"index;not null"`
	Temperature float64   `json:"temperature" gorm:"type:numeric(5,2)"`
	Humidity    float64   `json:"humidity" gorm:"type:numeric(5,2)"`
	Pressure    float64   `json:"pressure" gorm:"type:numeric(6,1)"`
	Topic       string    `json:"topic"`
	CreatedAt   time.Time `json:"created_at" gorm:"index"`
}

type SensorPayload struct {
	DeviceID    string  `json:"device_id"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Pressure    float64 `json:"pressure"`
}
