package mqtt

import (
	"encoding/json"
	"fmt"
	"github.com/detecc/detecctor-v2/model/configuration"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"time"
)

type (
	Topic string

	MessageHandler func(client Client, topicIds []string, payloadId uint16, payload interface{}, err error)

	// Client is an interface wrapper for a simple MQTT client.
	Client interface {
		Disconnect()
		Publish(topic Topic, message interface{}) error
		Subscribe(topic Topic, handler MessageHandler)
	}

	// ClientImpl concrete implementation of the Client, which is essentially a wrapper over the mqtt lib.
	ClientImpl struct {
		mqttClient mqtt.Client
	}
)

// NewMqttClient creates a wrapped mqtt Client with specific settings.
func NewMqttClient(clientSettings configuration.MqttBroker) Client {
	log.Info("Creating a new MQTT client..")
	broker := fmt.Sprintf("%s:%d", clientSettings.Host, clientSettings.Port)

	// Basic client settings
	ops := mqtt.NewClientOptions()
	ops.AddBroker(broker)
	ops.SetClientID(clientSettings.ClientId)
	ops.SetUsername(clientSettings.Username)
	ops.SetPassword(clientSettings.Password)

	// Connection settings
	ops.SetKeepAlive(30 * time.Second)
	ops.SetAutoReconnect(true)

	ops.SetOnConnectHandler(func(client mqtt.Client) {
		log.Info("Connected to broker")
	})

	ops.SetDefaultPublishHandler(func(client mqtt.Client, message mqtt.Message) {
		log.Printf("Received message %s from topic %s", message.Payload(), message.Topic())
	})

	ops.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		log.Info("Disconnected from broker:", err)
	})

	// Connect to the MQTT broker
	client := mqtt.NewClient(ops)
	client.Connect().Wait()

	return &ClientImpl{
		mqttClient: client,
	}
}

func (c *ClientImpl) Disconnect() {
	log.Debug("Disconnecting the MQTT client")
	c.mqttClient.Disconnect(100)
}

// Publish a new message to a topic
func (c *ClientImpl) Publish(topic Topic, message interface{}) error {
	logInfo := log.WithFields(log.Fields{
		"topic":   topic,
		"message": message,
	})
	logInfo.Debug("Publishing a message to topic")

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		logInfo.Errorf("Error marshalling the message: %v", err)
		return err
	}

	token := c.mqttClient.Publish(string(topic), 1, false, jsonMessage)
	go func(token mqtt.Token) {
		if token.Error() != nil {
			log.Warn(token.Error())
		}
	}(token)
	return nil
}

// Subscribe to a topic
func (c *ClientImpl) Subscribe(topic Topic, handler MessageHandler) {
	logInfo := log.WithField("topic", topic)
	logInfo.Debug("Subscribing to a topic")

	token := c.mqttClient.Subscribe(string(topic), 1, func(client mqtt.Client, message mqtt.Message) {
		var data interface{}

		// Transform the payload to the object and pass it to the handler function for further processing
		err := json.Unmarshal(message.Payload(), &data)
		if err != nil {
			logInfo.Errorf("Error parsing the data: %v", err)
			return
		}

		// Parse the topic and get the Ids based on the original topic.
		ids, err := GetIdsFromTopic(message.Topic(), topic)
		if err != nil {
			logInfo.Errorf("Error getting the topic info: %v", err)
			return
		}

		handler(c, ids, message.MessageID(), data, err)
	})

	go func(token mqtt.Token) {
		token.Wait()
		if token.Error() != nil {
			log.Warn(token.Error())
		}
	}(token)
}
