package notification

type (
	ProxyMessage struct {
		ChatId    string
		Username  string
		Message   string
		MessageId string
	}

	MessageBuilder struct {
		buildArgs []messageHandler
	}

	messageHandler func(*ProxyMessage)
)

func NewMessageBuilder() *MessageBuilder {
	return &MessageBuilder{}
}

func (b *MessageBuilder) WithMessage(messageId, message string) *MessageBuilder {
	b.buildArgs = append(b.buildArgs, func(msg *ProxyMessage) {
		msg.Message = message
		msg.MessageId = messageId
	})
	return b
}

func (b *MessageBuilder) WithId(chatId string) *MessageBuilder {
	b.buildArgs = append(b.buildArgs, func(msg *ProxyMessage) {
		msg.ChatId = chatId
	})
	return b
}

func (b *MessageBuilder) FromUser(username string) *MessageBuilder {
	b.buildArgs = append(b.buildArgs, func(msg *ProxyMessage) {
		msg.Username = username
	})
	return b
}

//Build builds the Reply object
func (b *MessageBuilder) Build() ProxyMessage {
	emp := ProxyMessage{
		ChatId:    "",
		Username:  "",
		Message:   "",
		MessageId: "",
	}

	for _, a := range b.buildArgs {
		a(&emp)
	}

	defer func() {
		// clear any actions
		b.buildArgs = []messageHandler{}
	}()

	return emp
}
