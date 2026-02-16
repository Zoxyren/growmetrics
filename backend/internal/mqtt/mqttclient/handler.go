package mqttclient

import (
	"encoding/json"
	"fmt"
	"log/slog"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"esp32/backend/internal/models"
)

func (a *MQTTAdapter) RecieveTopics(topic string, qos byte) mqtt.Token {
	token := a.client.Subscribe(topic, qos, a.RecieveMessage)
	token.Wait()
	fmt.Printf("Subscription to the topic %s successfully\n", topic)
	return token
}

func (a *MQTTAdapter) RecieveMessage(c mqtt.Client, msg mqtt.Message) {
	var payload models.SensorPayload

	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		slog.Error("JSON Unmarshal fehlgeschlagen", "error", err)
		return
	}

	err := a.sensorService.ProcessIncomingMetric(payload, msg.Topic())

	if err != nil {
		slog.Error("Verarbeitung fehlgeschlagen", "error", err)
	} else {
		fmt.Printf("Erfolgreich verarbeitet: %s\n", payload.DeviceID)
	}

	msg.Ack()
}

