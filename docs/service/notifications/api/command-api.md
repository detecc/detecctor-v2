# 📜 The guide to commands

## Command struct

The `Command` struct is created from a message that is sent to the bot and must have a `/` prefix. The message is split
by spaces; the first element in the split array is `Command.Name`, while the other elements of the array represent
the `Command.Args` (arguments).

The arguments are passed to the plugin which handles the `Command`. The `ChatId` is used to map the plugin response to
the chat.

```go
package command

type Command struct {
	Name   string
	Args   []string
	ChatId string
}
```

## 🏗️ Building the commands with CommandBuilder

To easily construct the commands from the bot, use the `CommandBuilder`.

```go
package example

import "github.com/detecc/detecctor-v2/model/command"

func CreateCommand() {
	builder := command.NewCommandBuilder()
	replyMessage := builder.WithId("chatId1").FromUser("userName").WithMessage("messageId", "messageText").Build()
}
```