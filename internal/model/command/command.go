package command

const (
	AuthShortCommand        = "auth"
	AuthCommand             = "authorize"
	DeAuthCommand           = "deauth"
	SubscribeCommand        = "subscribe"
	SubscribeShortCommand   = "sub"
	UnsubscribeCommand      = "unsubscribe"
	UnsubscribeShortCommand = "unsub"
)

type (
	// Command consists of a Name and Args.
	// Example of a command: "/get_status serviceNode1 serviceNode2".
	// The command name is "/get_status", the arguments are ["serviceNode1", "serviceNode2"].
	Command struct {
		Name      string   `json:"name" validate:"required"`
		Args      []string `json:"args" validate:"required"`
		MessageId string   `json:"messageId" validate:"required"`
		ChatId    string   `json:"chatId" validate:"required"`
	}
)
