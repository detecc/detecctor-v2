package proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/detecc/detecctor-v2/database"
	"github.com/detecc/detecctor-v2/internal/mqtt"
	"github.com/detecc/detecctor-v2/model/command"
	"github.com/detecc/detecctor-v2/model/command/logs"
	"github.com/detecc/detecctor-v2/model/reply"
	"github.com/detecc/detecctor-v2/service/notification"
	"github.com/detecc/detecctor-v2/service/notification/bot"
	mqtt2 "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
)

const (
	ReceiveNotificationTopic = mqtt.Topic("chat/+/notify")
)

var once = sync.Once{}
var proxy *Proxy

type (
	// Proxy acts as a proxy between the Bot implementation and the notification service.
	// It performs all the database operations and handles the bot communication.
	Proxy struct {
		bot          bot.Bot
		client       mqtt.Client
		replyBuilder *reply.Builder
	}
)

//NewProxy creates a new proxy for a bot and
func NewProxy(bot bot.Bot, client mqtt.Client) *Proxy {
	once.Do(func() {
		proxy = &Proxy{
			bot:          bot,
			client:       client,
			replyBuilder: reply.NewReplyBuilder(),
		}
	})
	return proxy
}

func GetProxy() *Proxy {
	return proxy
}

// Start the bot and listening for messages from the bot.
func (proxy *Proxy) Start(ctx context.Context) {
	log.Info("Starting the proxy...")
	defer log.Info("Stopping the proxy..")

	handler := func(client mqtt2.Client, message mqtt2.Message) {
		replyMsg := reply.Reply{}
		err := json.Unmarshal(message.Payload(), &replyMsg)
		if err != nil {
			log.WithFields(log.Fields{
				"topic":   message.Topic(),
				"payload": message.Payload(),
			}).Errorln("Cannot convert the payload to reply")
			return
		}

		// if the message is translatable, translate it before replying to the bot
		if replyMsg.ReplyType == reply.TranslatableMessage {
			replyMessage, err := TranslateReplyMessage(replyMsg.ChatId, replyMsg.Content)
			if err != nil {
				log.WithFields(log.Fields{
					"chatId":  replyMsg.ChatId,
					"content": replyMsg.Content,
				}).Error("Cannot translate the reply message")
				log.Println(err)
				return
			}

			replyMsg.ReplyType = reply.TypeMessage
			replyMsg.Content = replyMessage
		}

		proxy.bot.ReplyToChat(replyMsg)
	}
	proxy.client.Subscribe(ReceiveNotificationTopic, handler)

	proxy.bot.Start()
	go proxy.listenForMessages(ctx)
	proxy.bot.ListenToChannels()
}

// cleanup cleans up the bot and other components
func (proxy *Proxy) cleanup() {
	log.Debug("Cleaning up the proxy & bot")
}

// listenForMessages Listens and to messages from the Bot implementation and processes the message.
func (proxy *Proxy) listenForMessages(ctx context.Context) {
	log.Info("Starting to listen for messages..")
	for {
		select {
		case message := <-proxy.bot.GetMessageChannel():
			// ignore any non-Message Updates
			proxy.processMessage(message)
			break
		case <-ctx.Done():
			log.Info("Stopping the proxy listener..")
			proxy.cleanup()
		}
	}
}

// processMessage processes an incoming bot message.
// Adds the necessary information to the database and forwards the command to an appropriate service.
// If an error occurs during processing, the proxy will send a response with the error to the bot.
func (proxy *Proxy) processMessage(message notification.ProxyMessage) {
	log.WithField("message", message).Debug("Processing a message")

	err := database.AddChatIfDoesntExist(message.ChatId, message.Username)
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"chatId": message.ChatId,
		}).Error("Error adding a chat")
	}

	// Add a new message to database
	_, err = database.NewMessage(message.ChatId, message.MessageId, message.Message)
	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"chatId":    message.ChatId,
			"messageId": message.MessageId,
		}).Error("Error adding a message")
		return
	}

	// Update last message id
	err = database.UpdateLastMessageId(message.MessageId)
	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"messageId": message.MessageId,
		}).Error("Error updating the lastMessageId")
		return
	}

	cmd, err := parseCommand(message.Message, message.ChatId, message.MessageId)
	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"messageId": message.MessageId,
			"message":   message.Message,
		}).Error("Error parsing a message as a command")

		// If the command is invalid, notify the user
		replyMessage := proxy.replyBuilder.TypeMessage().ForChat(message.ChatId).WithContent(fmt.Sprintf("%s is not a command", message)).Build()
		proxy.bot.ReplyToChat(replyMessage)
		return
	}

	// Check if command has / prefix and remove it because of mqtt topics
	pluginName := cmd.Name
	if strings.HasPrefix(pluginName, "/") {
		pluginName = strings.TrimPrefix(pluginName, "/")
	}

	// Default publish topic
	publishTopic := fmt.Sprintf("plugin/%s/execute", pluginName)

	switch pluginName {
	case command.AuthCommand, command.AuthShortCommand:
		publishTopic = fmt.Sprintf("chat/%s/auth", cmd.ChatId)
		return
	case command.SubscribeShortCommand, command.SubscribeCommand:
		publishTopic = fmt.Sprintf("chat/%s/subscribe", cmd.ChatId)
		break
	case command.UnsubscribeShortCommand, command.UnsubscribeCommand:
		publishTopic = fmt.Sprintf("chat/%s/unsubscribe", cmd.ChatId)
		break
	}

	// Send the command to the designated topic that will handle the action
	err = proxy.client.Publish(mqtt.Topic(publishTopic), cmd)
	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"chatId":    message.ChatId,
			"messageId": message.MessageId,
			"topic":     publishTopic,
		}).Error("Error sending the command to the topic")
	}

	database.AddCommandLog(cmd, logs.WithErrors(err))
}
