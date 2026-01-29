package models

type SensorPayload struct {
	Device      string  `json:"device"`
	Humidity    float64 `json:"humidity"`
	Temperature float64 `json:"temperature"`
	Pressure    float64 `json:"pressure"`
	Topic       string  `json:"topic"`
}
