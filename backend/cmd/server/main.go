package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"esp32/backend/internal/core/domain"
	"esp32/backend/internal/database"
	"esp32/backend/internal/mqtt/mqttclient"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Info("Keine .env Datei gefunden, nutze System-Umgebungsvariablen")
	}
	db, err := database.NewPostgresConnection()
	if err != nil {
		slog.Error("failed to establish connection", "error:", err)
		os.Exit(1)
		return
	}
	defer db.Close()

	dbadapter := database.NewDatabaseAdapter(db)
	cacheAdapter := database.NewRedisAdapter(os.Getenv("REDIS_ADDR"))
	sensorService := domain.NewSensorService(dbadapter, cacheAdapter, cacheAdapter)

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
		metrics, err := dbadapter.GetAllMetrics()
		if err != nil {
			slog.Error("Failed to fetch Metrics", "error", err)
			return c.JSON(500, err.Error())
		}
		return c.JSON(200, metrics)
	})
	e.GET("/metrics/:id", func(c echo.Context) error {
		id := c.Param("id")
		data, err := sensorService.GetLatestData(c.Request().Context(), id)
		if err != nil {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusOK, data)
	})
	e.GET("/active-devices", func(c echo.Context) error {
		devices, err := sensorService.GetActiveDevices(c.Request().Context())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, devices)
	})
	fmt.Println("Backend aktiv. Webserver startet auf :1323...")
	if err := e.Start(":1323"); err != nil {
		slog.Error("Webserver konnte nicht starten", "error", err)
	}
}
