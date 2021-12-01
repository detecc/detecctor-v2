package bot

import (
	"context"
	"github.com/detecc/detecctor-v2/model/configuration"
	"github.com/detecc/detecctor-v2/model/reply"
	"github.com/detecc/detecctor-v2/service/notification"
	t "github.com/detecc/detecctor-v2/service/notification/providers/telegram"
	log "github.com/sirupsen/logrus"
)

const (
	TelegramBot = "telegram"
)

type (
	// Bot represents a chatbot (e.g. Telegram, Discord, Slack bot).
	Bot interface {
		// Start should initialize the bot to listen to the chat.
		Start()
		// ListenToChannels is called after the start function. It monitors the chat and should send the message data to the Message channel.
		ListenToChannels(ctx context.Context)
		// ReplyToChat after receiving the command results or if an error occurs.
		ReplyToChat(replyMessage reply.Reply)
		// GetMessageChannel returns the channel the bot is supposed to notify the proxy of the incoming message.
		GetMessageChannel() <-chan notification.ProxyMessage
	}
)

// NewBot create a new telegram bot.
func NewBot(botConfiguration configuration.BotConfiguration) Bot {
	log.Debugf("Creating a new bot with type: %s", botConfiguration.Type)

	switch botConfiguration.Type {
	case TelegramBot:
		return &t.Telegram{
			Token: botConfiguration.Token,
		}
	default:
		log.Fatalf("Unsupported bot type: %s", botConfiguration.Type)
		return nil
	}
}
