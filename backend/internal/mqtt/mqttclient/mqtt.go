package mqttclient

import (
	"fmt"
	"log/slog"

	// Importiere das core-Package für das SensorService Interface
	"esp32/backend/internal/core"

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
	client        mqtt.Client
	config        MQTTBroker
	sensorService core.SensorService
}

func NewMQTTAdapter(cfg MQTTBroker, svc core.SensorService) *MQTTAdapter {
	return &MQTTAdapter{
		config:        cfg,
		sensorService: svc,
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
		slog.Info("Verbunden mit MQTT-Broker", "url", a.config.MQTTBroker)
	}

	a.client = mqtt.NewClient(opts)
	if token := a.client.Connect(); token.Wait() && token.Error() != nil {
		slog.Error(
			"cannot establish connection to MQTT broker",
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

