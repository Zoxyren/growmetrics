package mqttclient

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTBroker struct {
	MQTTBroker string
	MQTTUser   string
	MQTTPW     string
	MQTTTopic  string
	ClientID   string
}

// Todo: .env variables instead of clear
// MQTTBroker{
// QTTBroker: "8b758ea22b9f4c0f94ac43c9b09a254f.s1.eu.hivemq.cloud",
//MQTTUser:   "esp32_1",
//MQTTPW:     "Testuser123",
//MQTTTopic:  "test",
//}

func EstablishESPConnection(config MQTTBroker) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.MQTTBroker)
	opts.SetClientID(config.ClientID)
	opts.SetUsername(config.MQTTUser)
	opts.SetPassword(config.MQTTPW)

	opts.OnConnect = func(c mqtt.Client) {
		fmt.Println("✅ Verbunden mit MQTT-Broker:", config.MQTTBroker)
	}
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return client
}
