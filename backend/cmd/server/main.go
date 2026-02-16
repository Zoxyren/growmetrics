package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"esp32/backend/internal/core/domain"
	"esp32/backend/internal/database"
	"esp32/backend/internal/mqtt/mqttclient"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Fehler beim Laden der .env Datei")
	}

	db, err := database.NewPostgresConnection()
	if err != nil {
		slog.Error("failed to establish connection", "error:", err)
		os.Exit(1)
		return
	}
	defer db.Close()

	dbadapter := database.NewDatabaseAdapter(db)

	sensorService := domain.NewSensorService(dbadapter)

	cfg := mqttclient.MQTTBroker{
		MQTTBroker: os.Getenv("MQTT_BROKER"),
		MQTTUser:   os.Getenv("MQTT_USER"),
		MQTTPW:     os.Getenv("MQTT_PW"),
		MQTTTopic:  os.Getenv("MQTT_TOPIC"),
		MQTTPort:   8883,
		ClientID:   os.Getenv("MQTT_CLIENT_ID"),
	}

	adapter := mqttclient.NewMQTTAdapter(cfg, sensorService)
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

	e.GET("/data", func(c echo.Context) error {
		// Auch hier nutzt du den Adapter, der das Interface erfüllt
		metrics, err := dbadapter.GetAllMetrics()
		if err != nil {
			slog.Error("Failed to fetch Metrics", "error", err)
			return c.JSON(500, err.Error())
		}
		return c.JSON(200, metrics)
	})

	go func() {
		if err := e.Start(":1323"); err != nil {
			fmt.Printf("Webserver Fehler: %v\n", err)
		}
	}()

	fmt.Println("Backend aktiv. Sende jetzt Daten von Wokwi.")
	select {}
}
