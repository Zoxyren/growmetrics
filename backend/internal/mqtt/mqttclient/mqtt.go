package mqttclient

import (
	"fmt"
	"log/slog"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTBroker struct {
	MQTTBroker string
	MQTTUser   string
	MQTTPW     string
	MQTTTopic  string
	ClientID   string
	MQTTPort   int
}

// Todo: .env variables instead of clear
var config = MQTTBroker{
	MQTTBroker: "8b758ea22b9f4c0f94ac43c9b09a254f.s1.eu.hivemq.cloud",
	MQTTUser:   "oliver3",
	MQTTPW:     "tESTUSER1234",
	MQTTTopic:  "esp32/oliver1/metrics",
	MQTTPort:   8883,
}

func EstablishESPConnection(config MQTTBroker) mqtt.Client {
	opts := mqtt.NewClientOptions()
	// opts.AddBroker(config.MQTTBroker)
	brokerURL := fmt.Sprintf("tls://%s:%d", config.MQTTBroker, config.MQTTPort)
	opts.AddBroker(brokerURL)
	opts.SetClientID(config.ClientID)
	opts.SetUsername(config.MQTTUser)
	opts.SetPassword(config.MQTTPW)

	opts.OnConnect = func(c mqtt.Client) {
		fmt.Println("✅ Verbunden mit MQTT-Broker:", config.MQTTBroker)
	}
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		slog.Error(
			"cannot establish connection to ESP",
			slog.Any("reason", token.Error()),
			slog.String("broker_url", brokerURL),
		)
		panic(error.Error)
	}
	return client
}
