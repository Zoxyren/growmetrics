package main

import (
	"fmt"

	"esp32/backend/mqttclient"
)

func main() {
	config := mqttclient.MQTTBroker{
		MQTTBroker: "8b758ea22b9f4c0f94ac43c9b09a254f.s1.eu.hivemq.cloud",
		MQTTUser:   "oliver1",
		MQTTPW:     "tESTUSER1234",
		MQTTTopic:  "esp32/oliver1/metrics",
		MQTTPort:   8883,
		ClientID:   "GoClient-12345",
	}

	client := mqttclient.EstablishESPConnection(config)
	defer client.Disconnect(250)

	token := client.Publish(config.MQTTTopic, 0, false, "Hallo MQTT!")
	token.Wait()

	fmt.Println("Programm läuft und hält die Verbindung. Beenden mit Ctrl+C.")
	select {}
}
