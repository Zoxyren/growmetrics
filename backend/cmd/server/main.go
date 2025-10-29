package main

import (
	"fmt"
	"log/slog"

	"github.com/labstack/echo/v4"

	"esp32/backend/internal/mqtt/mqttclient"
)

func main() {
	go func() {
		config := mqttclient.MQTTBroker{
			MQTTBroker: "8b758ea22b9f4c0f94ac43c9b09a254f.s1.eu.hivemq.cloud",
			MQTTUser:   "oliver3",
			MQTTPW:     "tESTUSER1234",
			MQTTTopic:  "esp32/oliver1/metrics",
			MQTTPort:   8883,
			ClientID:   "GoClient-12345",
		}

		client := mqttclient.EstablishESPConnection(config)
		defer client.Disconnect(250)
		subscribeToken := mqttclient.RecieveTopics(client, "esp32/oliver1/metrics", byte(1))

		if subscribeToken.Wait() && subscribeToken.Error() != nil {
			slog.Error("Failed to subscribe to topic!", "error", subscribeToken.Error)
		}

		fmt.Println("Programm läuft und hält die Verbindung. Beenden mit Ctrl+C.")
		print(subscribeToken)
		select {}
	}()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(202, "request sucesfully")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
