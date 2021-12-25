package reply

// constants for the Reply.ReplyType
const (
	TypeMessage = iota
	TranslatableMessage
	TypePhoto
	TypeAudio
)

type (
	// Reply is a struct used to parse results to the ReplyChannel in Bot.
	Reply struct {
		// Each reply must contain a ChatId - a chat to reply to.
		ChatId string `json:"chatId"`
		// The ReplyType must be a constant defined in the package.
		ReplyType int `json:"replyType"`
		// Content must be cast after determining the type to send to Bot.
		Content interface{} `json:"content,omitempty"`
	}
)
