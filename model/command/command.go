package command

import (
	"fmt"
	"strings"
)

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

	Builder struct {
		buildActions []handler
	}

	handler func(cmd *Command)
)

//NewCommandBuilder - constructor
func NewCommandBuilder() *Builder {
	return &Builder{
		buildActions: []handler{},
	}
}

//WithName sets name of the command
func (b *Builder) WithName(value string) *Builder {
	b.buildActions = append(b.buildActions, func(cmd *Command) {
		if !strings.HasPrefix(value, "/") {
			value = fmt.Sprintf("/%s", value)
		}

		cmd.Name = value
	})
	return b
}

//WithArgs sets arguments of the command
func (b *Builder) WithArgs(args []string) *Builder {
	b.buildActions = append(b.buildActions, func(cmd *Command) {
		if args == nil {
			return
		}
		cmd.Args = args
	})
	return b
}

//FromChat sets chatId
func (b *Builder) FromChat(chatId string) *Builder {
	b.buildActions = append(b.buildActions, func(cmd *Command) {
		cmd.ChatId = chatId
	})
	return b
}

//Id sets messageId
func (b *Builder) Id(messageId string) *Builder {
	b.buildActions = append(b.buildActions, func(cmd *Command) {
		cmd.MessageId = messageId
	})
	return b
}

//Build builds the Command object
func (b *Builder) Build() Command {
	emp := Command{
		Name:   "",
		Args:   []string{},
		ChatId: "0",
	}

	for _, a := range b.buildActions {
		a(&emp)
	}

	defer func() {
		// clear any previous actions
		b.buildActions = []handler{}
	}()

	return emp
}
