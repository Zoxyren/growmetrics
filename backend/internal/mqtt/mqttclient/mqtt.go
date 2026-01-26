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

type MQTTAdapter struct {
	client mqtt.Client
	config MQTTBroker
}

// Todo: .env variables instead of clear
var config = MQTTBroker{
	MQTTBroker: "8b758ea22b9f4c0f94ac43c9b09a254f.s1.eu.hivemq.cloud",
	MQTTUser:   "oliver3",
	MQTTPW:     "tESTUSER1234",
	MQTTTopic:  "esp32/oliver1/metrics",
	MQTTPort:   8883,
	ClientID:   "GoClient-12345",
}

func NewAdapter(cfg MQTTBroker) *MQTTAdapter {
	return &MQTTAdapter{
		config: cfg,
	}
}

func (a *MQTTAdapter) Connect() error {
	opts := mqtt.NewClientOptions()
	// opts.AddBroker(config.MQTTBroker)
	brokerURL := fmt.Sprintf("tls://%s:%d", a.config.MQTTBroker, a.config.MQTTPort)
	opts.AddBroker(brokerURL)
	opts.SetClientID(a.config.ClientID)
	opts.SetUsername(a.config.MQTTUser)
	opts.SetPassword(a.config.MQTTPW)

	opts.OnConnect = func(c mqtt.Client) {
		fmt.Println("✅ Verbunden mit MQTT-Broker:", config.MQTTBroker)
	}
	a.client = mqtt.NewClient(opts)
	if token := a.client.Connect(); token.Wait() && token.Error() != nil {
		slog.Error(
			"cannot establish connection to ESP",
			slog.Any("reason", token.Error()),
			slog.String("broker_url", brokerURL),
		)
		return token.Error()
	}
	return nil
}

func (a *MQTTAdapter) Disconnect(quiesce uint) {
	if a.client != nil && a.client.IsConnected() {
		a.client.Disconnect(quiesce)
		fmt.Println("MQTT Verbindung erfolgreich geschlossen")
	}
}
