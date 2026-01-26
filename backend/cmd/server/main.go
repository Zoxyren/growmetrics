package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"

	_ "github.com/joho/godotenv/autoload"

	"esp32/backend/internal/mqtt/mqttclient"
)

func main() {
	cfg := mqttclient.MQTTBroker{
		MQTTBroker: os.Getenv("MQTTBroker"),
		MQTTUser:   os.Getenv("MQTTUser"),
		MQTTPW:     os.Getenv("MQTTPW"),
		MQTTTopic:  os.Getenv("MQTTTopic"),
		MQTTPort:   8883,
		ClientID:   os.Getenv("ClientID"),
	}
	adapter := mqttclient.NewAdapter(cfg)
	if err := adapter.Connect(); err != nil {
		slog.Error("Konnte MQTT nicht verbinden", "error:", err)
	}
	go func() {
		defer adapter.Disconnect(250)
		subscribeToken := adapter.RecieveTopics("esp32/oliver1/metrics", 1)

		if subscribeToken.Wait() && subscribeToken.Error() != nil {
			slog.Error("Failed to subscribe to topic!", "error", subscribeToken.Error)
		}

		fmt.Println("Programm läuft und hält die Verbindung. Beenden mit Ctrl+C.")
		select {}
	}()
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(202, "request sucesfully")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
