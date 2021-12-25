package i18n

type (
	// TranslationMap is a wrapper to easily use i18n in the server.
	TranslationMap struct {
		MessageId string
		Plural    interface{}
		Data      map[string]interface{}
	}

	// TranslationOptions for TranslationMap
	TranslationOptions func(*TranslationMap)
)

// NewTranslationMap makes a map that can be interpreted in the SendMessageToChat's content argument as a message to be translated.
func NewTranslationMap(messageId string, opts ...TranslationOptions) TranslationMap {
	translatedMap := &TranslationMap{
		MessageId: messageId,
		Plural:    nil,
		Data:      make(map[string]interface{}),
	}

	for _, opt := range opts {
		opt(translatedMap)
	}

	return *translatedMap
}

// WithPlural add plurality
func WithPlural(plural interface{}) TranslationOptions {
	return func(translationMap *TranslationMap) {
		translationMap.Plural = plural
	}
}

// AddData Add key-value pairs to the data map
func AddData(key string, value interface{}) TranslationOptions {
	return func(translationMap *TranslationMap) {
		translationMap.Data[key] = value
	}
}
