package i18n

import (
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	translationsPath = "./i18n/translations"
)

var (
	bundle             *i18n.Bundle
	defaultMessages    map[string]i18n.Message
	Matcher            language.Matcher
	supportedLanguages = []language.Tag{
		language.English, // The first language is used as fallback.
	}
)

func init() {
	once := sync.Once{}
	once.Do(func() {
		defaultMessages = make(map[string]i18n.Message)

		bundle = i18n.NewBundle(language.English)
		bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

		loadTranslations()

		// command and cmd messages
		AddDefaultMessage(i18n.Message{
			ID:    "UnsupportedCommand",
			Other: "Command {{.Command}} is unsupported.",
		})
		AddDefaultMessage(i18n.Message{
			ID:    "InvalidArguments",
			One:   "Invalid argument.",
			Other: "Invalid arguments.",
		})
		AddDefaultMessage(i18n.Message{
			ID:    "PluginExecutionFailed",
			Other: "An error occurred during the command execution: {{.Error}}",
		})
		AddDefaultMessage(i18n.Message{
			ID:    "InvalidPluginType",
			Other: "Plugin {{.Plugin}} has an invalid cmd type: {{.PluginType}}",
		})
		// chat authorization messages
		AddDefaultMessage(i18n.Message{
			ID:    "ChatUnauthorized",
			Other: "You are not authorized to send this command.",
		})
		AddDefaultMessage(i18n.Message{
			ID:    "ChatAuthorized",
			Other: "Chat successfully authorized.",
		})
		AddDefaultMessage(i18n.Message{
			ID:    "InvalidToken",
			Other: "Invalid or expired token.",
		})
		AddDefaultMessage(i18n.Message{
			ID:    "AuthorizationError",
			Other: "Error authorizing the chat.",
		})
		AddDefaultMessage(i18n.Message{
			ID:    "GeneratedToken",
			Other: "Generated a token for authentication.",
		})
		// client messages
		AddDefaultMessage(i18n.Message{
			ID:    "ClientDisconnected",
			Other: "Client {{.ServiceNodeKey}} has been disconnected at {{.Time}}.",
		})
		AddDefaultMessage(i18n.Message{
			ID:    "UnableToSendMessage",
			Other: "Unable to send the message to {{.ServiceNodeKey}}: {{.Error}}",
		})
		// subscription messages
		AddDefaultMessage(i18n.Message{
			ID:    "SubscriptionSuccess",
			One:   "Successfully subscribed to command and/or node.",
			Other: "Successfully subscribed to commands and/or nodes.",
		})
		AddDefaultMessage(i18n.Message{
			ID:    "UnsubscribeSuccess",
			One:   "Successfully unsubscribed from command and/or node.",
			Other: "Successfully unsubscribed from commands and/or nodes.",
		})

		AddDefaultMessage(i18n.Message{
			ID:    "SubscriptionFail",
			Other: "An error occurred during subscription: {{.Error}}",
		})
		AddDefaultMessage(i18n.Message{
			ID:    "UnsubscribeFail",
			Other: "An error occurred during unsubscribing: {{.Error}}",
		})
	})
}

// loadTranslations loads all available translations from the translations folder into the bundle.
func loadTranslations() {
	log.Info("Loading translations..")
	defer log.Info("Successfully loaded translations:", supportedLanguages)

	err := filepath.Walk(translationsPath, func(path string, info os.FileInfo, err error) error {
		// Load all active.*.yaml translations into the bundle
		if info != nil && !info.IsDir() && strings.Contains(info.Name(), "active.") {
			return loadTranslation(path, info)
		}

		return nil
	})

	if err != nil {
		log.WithError(err).Errorf("An error occured when loading translations")
	}

	// Create a matcher based on imported translation files.
	Matcher = language.NewMatcher(supportedLanguages)
}

func loadTranslation(path string, info os.FileInfo) error {
	// Example file name: active.en.yaml -> get "en"
	strs := strings.Split(info.Name(), ".")
	if len(strs) < 2 {
		return fmt.Errorf("invalid file name")
	}

	// The language is second to last
	lang := strs[len(strs)-2]

	log.Debugf("Loading a translation: %s", lang)

	err := addLanguageSupport(lang)
	if err != nil {
		return err
	}

	// Load the translation file
	_, err = bundle.LoadMessageFile(path)
	return err
}

func addLanguageSupport(lang string) error {
	tag, err := language.Parse(lang)
	if err != nil {
		return err
	}

	supportedLanguages = append(supportedLanguages, tag)
	return nil
}

// AddDefaultMessage exposes the API for plugins
func AddDefaultMessage(message i18n.Message) {
	defaultMessages[message.ID] = message
}

// Localize translates the message based on the language of the chat.
func Localize(lang string, messageId string, data map[string]interface{}, plural interface{}) (string, error) {
	log.Debugf("Translating message %s to %s", messageId, lang)

	tag, _ := language.MatchStrings(Matcher, lang)
	locale := i18n.NewLocalizer(bundle, tag.String())
	defaultMessage, ok := defaultMessages[messageId]
	if !ok {
		return "", fmt.Errorf("default message not found")
	}

	msg, err := locale.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &defaultMessage,
		TemplateData:   data,
		PluralCount:    plural,
	})
	if err != nil {
		return "", err
	}

	return msg, nil
}
