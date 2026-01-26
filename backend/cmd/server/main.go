package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"

	"esp32/backend/internal/mqtt/mqttclient"
)

func main() {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatal("Fehler beim Laden der .env Datei")
	}
	fmt.Printf("Verbinde zu Broker: %s\n", os.Getenv("MQTT_BROKER"))
	cfg := mqttclient.MQTTBroker{
		MQTTBroker: os.Getenv("MQTT_BROKER"),
		MQTTUser:   os.Getenv("MQTT_USER"),
		MQTTPW:     os.Getenv("MQTT_PW"),
		MQTTTopic:  os.Getenv("MQTT_TOPIC"),
		MQTTPort:   8883,
		ClientID:   os.Getenv("CLIENT_ID"),
	}
	adapter := mqttclient.NewAdapter(cfg)
	if err := adapter.Connect(); err != nil {
		slog.Error("Konnte MQTT nicht verbinden", "error:", err)
	}
	go func() {
		defer adapter.Disconnect(250)
		subscribeToken := adapter.RecieveTopics(cfg.MQTTTopic, 1)

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
