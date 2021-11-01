package database

import (
	"context"
	"fmt"
	"github.com/agrison/go-commons-lang/stringUtils"
	. "github.com/detecc/detecctor-v2/model/chat"
	"github.com/kamva/mgm/v3"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/text/language"
)

func getChat(filter interface{}) (*Chat, error) {
	chat := &Chat{}
	chatCollection := mgm.Coll(chat)

	// Get the first doc of a collection using a filter
	err := chatCollection.First(filter, chat)
	if err != nil {
		log.Println("Error querying chats:", err)
		return nil, err
	}
	return chat, nil
}

func getChatWithCtx(ctx context.Context, filter interface{}) (*Chat, error) {
	chat := &Chat{}

	// Get the first doc of a collection using a filter
	err := mgm.Coll(chat).FirstWithCtx(ctx, filter, chat)
	if err != nil {
		return nil, err
	}

	return chat, nil
}

func updateChat(chat *Chat) error {
	chatCollection := mgm.Coll(&Chat{})
	return chatCollection.Update(chat)
}

func updateChatWithCtx(ctx context.Context, chat *Chat) error {
	return mgm.Coll(&Chat{}).UpdateWithCtx(ctx, chat)
}

func IsChatAuthorized(chatId string) bool {
	log.WithField("chatId", chatId).Debug("Checking if chat authorized")

	chat, err := GetChatWithId(chatId)
	if err != nil {
		log.Debug("Error authenticating the chat:", err)
		return false
	}

	return chat.IsAuthorized
}

func AuthorizeChat(chatId string) error {
	log.WithField("chatId", chatId).Debug("Authorizing chat")

	return mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		chat, err := getChatWithCtx(sc, bson.M{"chatId": chatId})
		if err != nil {
			return err
		}

		chat.IsAuthorized = true

		err = updateChatWithCtx(sc, chat)
		if err != nil {
			return err
		}

		return session.CommitTransaction(sc)
	})
}

func RevokeChatAuthorization(chatId string) error {
	log.WithField("chatId", chatId).Debug("Revoking chat authorization")

	return mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		chat, err := getChatWithCtx(sc, bson.M{"chatId": chatId})
		if err != nil {
			return err
		}

		chat.IsAuthorized = false

		err = updateChatWithCtx(sc, chat)
		if err != nil {
			return err
		}

		return session.CommitTransaction(sc)
	})
}

func GetChatWithId(chatId string) (*Chat, error) {
	log.WithField("chatId", chatId).Debug("Getting a chat")
	return getChat(bson.M{"chatId": chatId})
}

func GetChats() ([]Chat, error) {
	return getChats(bson.M{})
}

func getChats(filter interface{}) ([]Chat, error) {
	var (
		chat    = &Chat{}
		results []Chat
	)
	cursor, err := mgm.Coll(chat).Find(mgm.Ctx(), filter)
	if err = cursor.All(mgm.Ctx(), &results); err != nil {
		return nil, err
	}
	return results, nil
}

func addChat(ctx context.Context, chatId string, name string) error {
	log.WithField("chatId", chatId).Debug("Adding a new chat")

	if stringUtils.IsEmpty(chatId) {
		return fmt.Errorf("chat id is empty")
	}

	chat := &Chat{
		ChatId:        chatId,
		Name:          name,
		IsAuthorized:  false,
		Language:      language.English.String(),
		Subscriptions: []Subscription{},
	}

	return mgm.Coll(&Chat{}).CreateWithCtx(ctx, chat)
}

func AddChatIfDoesntExist(chatId string, name string) error {

	return mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		_, err := getChatWithCtx(sc, bson.M{"chatId": chatId})

		switch err {
		case mongo.ErrNoDocuments:
			err := addChat(sc, chatId, name)
			if err != nil {
				return err
			}

			return session.CommitTransaction(sc)
		default:
			return fmt.Errorf("chat already exists")
		}
	})
}
