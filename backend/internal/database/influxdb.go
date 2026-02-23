package database

import (
	"os"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func NewInfluxDB() influxdb2.Client {
	token := os.Getenv("INFLUXDB_TOKEN")
	url := "http://localhost:8086"
	client := influxdb2.NewClient(url, token)

	return client
}

