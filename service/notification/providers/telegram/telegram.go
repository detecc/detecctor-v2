package bot

import (
	"context"
	"fmt"
	"github.com/detecc/detecctor-v2/database"
	"github.com/detecc/detecctor-v2/model/reply"
	"github.com/detecc/detecctor-v2/service/notification"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"strconv"
)

// Telegram is a wrapper for the Telegram bot API.
type Telegram struct {
	Token          string
	botAPI         *telegram.BotAPI
	messageChannel chan notification.ProxyMessage
	messageBuilder *notification.MessageBuilder
}

// Start listening to the bot updates and the updates from the notification service
func (t *Telegram) Start() {
	telegramBot, err := telegram.NewBotAPI(t.Token)
	if err != nil {
		log.Panic(err)
	}

	t.botAPI = telegramBot
	t.messageChannel = make(chan notification.ProxyMessage)
	t.messageBuilder = notification.NewMessageBuilder()
}

//GetMessageChannel returns the channel for outgoing messages
func (t *Telegram) GetMessageChannel() <-chan notification.ProxyMessage {
	return t.messageChannel
}

// ListenToChannels listens for incoming data from telegram bot messages
func (t *Telegram) ListenToChannels(ctx context.Context) {
	log.Infof("Authorized on account %s", t.botAPI.Self.UserName)
	message, err := database.GetStatistics().GetStatistics(nil)

	lastMessageId := 0
	if err == nil {
		messageId, err := strconv.Atoi(message.LastMessageId)
		if err == nil {
			lastMessageId = messageId
		} else {
			log.Errorf("Error converting message id to int:%v", err)
		}
	}

	u := telegram.NewUpdate(lastMessageId)
	u.Timeout = 60

	updates, err := t.botAPI.GetUpdatesChan(u)
	if err != nil {
		log.Errorf("Error receiving update channel from Telegram: %v", err)
		return
	}

	for {
		select {
		case update := <-updates:
			if update.Message == nil || update.Message.Entities == nil || len(*update.Message.Entities) == 0 {
				return
			}

			for _, entity := range *update.Message.Entities {
				if entity.Type == "bot_command" {
					chatId := fmt.Sprintf("%d", update.Message.Chat.ID)
					messageId := fmt.Sprintf("%d", update.Message.MessageID)

					t.messageBuilder.WithId(chatId).FromUser(update.Message.Chat.UserName).WithMessage(messageId, update.Message.Text)
					t.messageChannel <- t.messageBuilder.Build()
				}
			}
			break
		case <-ctx.Done():
			log.Info("Stopping the Telegram bot..")
			return
		}
	}
}

//ReplyToChat replies to a telegram chat
func (t *Telegram) ReplyToChat(replyMessage reply.Reply) {
	var msg telegram.Chattable

	chatId, err := strconv.Atoi(replyMessage.ChatId)
	if err != nil {
		log.Errorf("Error converting ChatId to int:%v", err)
		return
	}

	switch replyMessage.ReplyType {
	case reply.TypeMessage:
		if replyMessage.Content != nil {
			msg = telegram.NewMessage(int64(chatId), replyMessage.Content.(string))
		}
		break
	case reply.TypePhoto:
		msg = telegram.NewPhotoUpload(int64(chatId), replyMessage.Content)
		break
	default:
		return
	}

	if msg != nil {
		log.WithFields(log.Fields{
			"chatId":    replyMessage.ChatId,
			"replyType": replyMessage.ReplyType,
		}).Debug("Replying to chat")

		_, err := t.botAPI.Send(msg)
		if err != nil {
			log.Errorf("Error sending the message to chat: %v", err)
		}
	}
}
