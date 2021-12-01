package mongo

import (
	"context"
	"fmt"
	. "github.com/detecc/detecctor-v2/model/chat"
	"github.com/kamva/mgm/v3"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageRepository struct{}

func NewMessageRepository() *MessageRepository {
	return &MessageRepository{}
}

func (m *MessageRepository) GetMessageFromChat(ctx context.Context, chatId int) (*Message, error) {
	log.WithField("chatId", chatId).Debug("Getting messages from chat")

	return getMessage(bson.M{"chatId": chatId})
}

func (m *MessageRepository) GetMessagesFromChat(ctx context.Context, chatId string) ([]Message, error) {
	var (
		msg               = &Message{}
		messageCollection = mgm.Coll(msg)
		results           []Message
	)

	// find all messages with the chatId
	cursor, err := messageCollection.Find(mgm.Ctx(), bson.D{{"chatId", chatId}})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(mgm.Ctx(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (m *MessageRepository) GetMessageWithId(ctx context.Context, messageId string) (*Message, error) {
	log.WithField("messageId", messageId).Debug("Getting a message with id")

	return getMessage(bson.M{"messageId": messageId})
}

func (m *MessageRepository) NewMessage(ctx context.Context, chatId string, messageId string, content string) (*Message, error) {
	log.WithFields(log.Fields{
		"chatId":    chatId,
		"messageId": messageId,
		"content":   content,
	}).Debug("Creating a message")

	message := &Message{
		ChatId:    chatId,
		Content:   content,
		MessageId: messageId,
	}

	err := mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		err := addMessageWithCtx(sc, message)
		if err != nil {
			return err
		}

		return session.CommitTransaction(sc)
	})
	if err != nil {
		return nil, err
	}

	return message, nil
}

func getMessage(filter interface{}) (*Message, error) {
	msg := &Message{}

	err := mgm.Coll(&Message{}).First(filter, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func getMessageWithCtx(ctx context.Context, filter interface{}) (*Message, error) {
	msg := &Message{}

	err := mgm.Coll(&Message{}).FirstWithCtx(ctx, filter, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func addMessageWithCtx(ctx context.Context, message *Message) error {
	if message == nil {
		return fmt.Errorf("message cannot be nil pointer")
	}

	_, err := getMessageWithCtx(ctx, bson.M{"chatId": message.ChatId})
	switch err {
	case nil:
		return fmt.Errorf("duplicate message found")
	case mongo.ErrNoDocuments:
		return mgm.Coll(&Message{}).CreateWithCtx(ctx, message)
	default:
		return err
	}
}
