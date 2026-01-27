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

func NewAdapter(cfg MQTTBroker) *MQTTAdapter {
	return &MQTTAdapter{
		config: cfg,
	}
}

func (a *MQTTAdapter) Connect() error {
	opts := mqtt.NewClientOptions()
	brokerURL := fmt.Sprintf("tls://%s:%d", a.config.MQTTBroker, a.config.MQTTPort)
	opts.AddBroker(brokerURL)
	opts.SetClientID(a.config.ClientID)
	opts.SetUsername(a.config.MQTTUser)
	opts.SetPassword(a.config.MQTTPW)

	opts.OnConnect = func(c mqtt.Client) {
		fmt.Println("Verbunden mit MQTT-Broker:", a.config.MQTTBroker)
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
