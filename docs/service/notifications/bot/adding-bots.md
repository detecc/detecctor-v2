# 🤖 Adding support for Bots

All chat services' bots must implement the `Bot` interface with proper functionality in order to work with the
Detecctor.

```go
package bot

import (
	reply "github.com/detecc/detecctor-v2/model/reply"
	"github.com/detecc/detecctor-v2/notifications"
)

// Bot represents a chatbot (e.g. Telegram, Discord, Slack bot).
type Bot interface {
	// Start should initialize the bot to listen to the chat.
	Start()
	// ListenToChannels is called after the start function. It monitors the chat and should send the message data to the Message channel.
	ListenToChannels()
	// ReplyToChat after executing the command or if an error occurs.
	ReplyToChat(replyMessage reply.Reply)
	// GetMessageChannel returns the channel the bot is supposed to notify the proxy of the incoming message.
	GetMessageChannel() chan notifications.ProxyMessage
}
```

## ♻️Contributing

To add bot support with your implementation, create a PR. Code changes should include:

1. New file and implementation of the `Bot` interface with the name of the bot
2. Adding a bot type constant in the [bot file](../../bot/bot.go) or in the newly created file
3. Modifying the supported bot table in [readme](../../Readme.md#-supported-bots)
4. Adding docs for fixes/development