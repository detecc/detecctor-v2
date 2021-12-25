package proxy

import (
	"context"
	"fmt"
	"github.com/detecc/detecctor-v2/database/repositories"
	commandBuilder "github.com/detecc/detecctor-v2/internal/command"
	"github.com/detecc/detecctor-v2/internal/model/command"
	i18n2 "github.com/detecc/detecctor-v2/pkg/i18n"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
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
	cmdBuilder := commandBuilder.NewCommandBuilder()

	if len(args) == 1 {
		cmdBuilder.WithName(args[0]).FromChat(chatId).Id(messageId)
	} else {
		cmdBuilder.WithName(args[0]).WithArgs(args[1:]).Id(messageId).FromChat(chatId)
	}

	return cmdBuilder.Build(), nil
}

// TranslateReplyMessage remaps the content and translates the message using i18n.
// The translation is dependent on the chat language.
func TranslateReplyMessage(chatRepository repositories.ChatRepository, chatId string, content interface{}) (string, error) {
	var (
		translationMap = content.(i18n2.TranslationMap)
		messageId      = translationMap.MessageId
		data           = translationMap.Data
		plural         = translationMap.Plural
		logInfo        = log.WithFields(log.Fields{
			"chatId":    chatId,
			"messageId": messageId,
		})
	)

	logInfo.Debug("Translating a reply message")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	// Get the preferred language for the chat
	lang, err2 := chatRepository.GetLanguage(ctx, chatId)
	if err2 != nil {
		cancel()
		return "", err2
	}

	cancel()

	// Check if a translation for the language is available
	localize, err := i18n2.Localize(lang, messageId, data, plural)
	if err != nil {
		logInfo.WithError(err).Warnf("Error localizing the translationMap")
		return "", err
	}

	return localize, nil
}
