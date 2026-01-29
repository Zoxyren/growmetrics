package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"esp32/backend/internal/database"
	"esp32/backend/internal/mqtt/mqttclient"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatal("Fehler beim Laden der .env Datei")
	}

	cfg := mqttclient.MQTTBroker{
		MQTTBroker: os.Getenv("MQTT_BROKER"),
		MQTTUser:   os.Getenv("MQTT_USER"),
		MQTTPW:     os.Getenv("MQTT_PW"),
		MQTTTopic:  os.Getenv("MQTT_TOPIC"),
		MQTTPort:   8883,
		ClientID:   os.Getenv("MQTT_CLIENT_ID"),
	}

	fmt.Printf("Verbinde zu Broker: %s\n", cfg.MQTTBroker)

	db, err := database.NewPostgresConnection()
	if err != nil {
		slog.Error("failed to establish connection", "error:", err)
	}
	defer db.Close()

	adapter := mqttclient.NewAdapter(cfg, db)
	if err := adapter.Connect(); err != nil {
		slog.Error("Konnte MQTT nicht verbinden", "error:", err)
		return
	}

	adapter.RecieveTopics("esp32/oliver1/metrics", 1)
	fmt.Println("Warte auf Topic: esp32/oliver1/metrics...")

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "request successfully")
	})

	go func() {
		if err := e.Start(":1323"); err != nil {
			fmt.Printf("Webserver Fehler: %v\n", err)
		}
	}()

	fmt.Println("Backend aktiv. Sende jetzt Daten von Wokwi.")
	select {}
}
