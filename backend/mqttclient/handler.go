package mqttclient

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func RecieveTopics(c mqtt.Client, topic string, qos byte) mqtt.Token {
	token := c.Subscribe(topic, qos, RevieveMessage)
	token.Wait()
	fmt.Println("Subscription to the topic succesfully", topic)
	return token
}

func RevieveMessage(c mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Nachricht auf Topic %s: %s\n", msg.Topic(), msg.Payload())
	msg.Ack()
}
