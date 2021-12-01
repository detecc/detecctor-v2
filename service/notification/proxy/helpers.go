package proxy

import (
	"fmt"
	"github.com/detecc/detecctor-v2/database"
	"github.com/detecc/detecctor-v2/internal/i18n"
	"github.com/detecc/detecctor-v2/model/command"
	log "github.com/sirupsen/logrus"
	"strings"
)

// parseCommand parses the text as a command, where the command is structured as /command arg1 arg2 arg3.
// returns a Command struct containing the name of the command and the arguments provided: ["/command", "arg1", "arg2", "arg3"]
func parseCommand(text, chatId, messageId string) (command.Command, error) {
	log.WithFields(log.Fields{
		"chatId":    chatId,
		"messageId": messageId,
		"content":   text,
	}).Debug("Parsing a command")

	if !strings.HasPrefix(text, "/") {
		return command.Command{}, fmt.Errorf("not a command: %s", text)
	}

	args := strings.Split(text, " ")
	cmdBuilder := command.NewCommandBuilder()

	if len(args) == 1 {
		cmdBuilder.WithName(args[0]).FromChat(chatId).Id(messageId)
	} else {
		cmdBuilder.WithName(args[0]).WithArgs(args[1:]).Id(messageId).FromChat(chatId)
	}

	return cmdBuilder.Build(), nil
}

// TranslateReplyMessage remaps the content and translates the message using i18n.
// The translation is dependent on the chat language.
func TranslateReplyMessage(chatId string, content interface{}) (string, error) {
	translationMap := content.(i18n.TranslationMap)
	messageId := translationMap.MessageId
	data := translationMap.Data
	plural := translationMap.Plural

	log.WithFields(log.Fields{
		"chatId":    chatId,
		"messageId": messageId,
	}).Debug("Translating a reply message")

	// Get the preferred language for the chat
	lang, err2 := database.GetChatRepository().GetLanguage(chatId)
	if err2 != nil {
		return "", err2
	}

	// Check if a translation for the language is available
	localize, err := i18n.Localize(lang, messageId, data, plural)
	if err != nil {
		log.Println("Error localizing the translationMap", err)
		return "", err
	}

	return localize, nil
}
