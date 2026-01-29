package mqttclient

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"esp32/backend/internal/models"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func (a *MQTTAdapter) RecieveTopics(topic string, qos byte) mqtt.Token {
	token := a.client.Subscribe(topic, qos, a.RecieveMessage)
	token.Wait()
	fmt.Println("Subscription to the topic succesfully", topic)
	return token
}

func (a *MQTTAdapter) RecieveMessage(c mqtt.Client, msg mqtt.Message) {
	var data models.SensorPayload
	fmt.Println("!!! NACHRICHT EMPFANGEN !!!")

	fmt.Printf("Raw Payload: %s\n", string(msg.Payload()))
	if err := json.Unmarshal(msg.Payload(), &data); err != nil {
		slog.Error("JSON Unmarshal fehlgeschlagen", "error", err)
		return
	}

	rawPayload := string(msg.Payload())
	topic := msg.Topic()
	// kein insert warum(?)
	query := `
		INSERT INTO sensor_data (device_id, temperature, humidity, pressure, topic)
		VALUES ($1, $2, $3, $4, $5)`

	res, err := a.db.Exec(query,
		data.Device,
		data.Temperature,
		data.Humidity,
		data.Pressure,
		topic,
	)

	if err != nil {
		slog.Error("DB-Insert fehlgeschlagen", "error", err)
	} else {
		rows, _ := res.RowsAffected()
		fmt.Printf("Gespeichert | Gerät: %s | Zeilen: %d | Payload: %s\n", data.Device, rows, rawPayload)
	}

	msg.Ack()
}

// Todo: Add Support for database-Insert, change table and columns for the correct data
