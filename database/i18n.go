package database

import (
	"fmt"
	"github.com/detecc/detecctor-v2/internal/i18n"
	"github.com/kamva/mgm/v3"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/text/language"
)

// GetLanguage for a chat
func GetLanguage(chatId string) (string, error) {
	log.WithField("chatId", chatId).Debug("Getting a language for chat")

	chat, err := getChat(bson.M{"chatId": chatId})
	if err != nil {
		return "", err
	}

	return chat.Language, nil
}

// SetLanguage changes the language preference from a default one.
func SetLanguage(chatId string, lang string) error {
	log.WithFields(log.Fields{
		"chatId":   chatId,
		"language": lang,
	}).Debug("Updating chat language")

	return mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		chat, err := getChatWithCtx(sc, bson.M{"chatId": chatId})
		if err != nil {
			return err
		}

		tag, _ := language.MatchStrings(i18n.Matcher, lang)
		if tag.String() == lang { // if the language is supported, update the chat
			chat.Language = tag.String()

			err = updateChatWithCtx(sc, chat)
			if err != nil {
				return err
			}

			return session.CommitTransaction(sc)
		}

		return fmt.Errorf("unsupported language: %v", lang)
	})
}
