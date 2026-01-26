package mqttclient

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func (a *MQTTAdapter) RecieveTopics(topic string, qos byte) mqtt.Token {
	token := a.client.Subscribe(topic, qos, a.RevieveMessage)
	token.Wait()
	fmt.Println("Subscription to the topic succesfully", topic)
	return token
}

func (a *MQTTAdapter) RevieveMessage(c mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Nachricht auf Topic %s: %s\n", msg.Topic(), msg.Payload())
	msg.Ack()
}
