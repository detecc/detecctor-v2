# The guide to bot replies

## Reply struct

The `Reply` struct is used in server-to-bot communication. The `ChatId` is used to reply to a specific chat.
The `ReplyType` represents one of the constants and specifies the type of `Content`. The `Content`
is a generic representation of the data, usually produced by the plugin (e.g. an Image, Text, Audio file, etc.)

```go
package reply

// constants for the Reply type
const (
	TypeMessage = 0
	TypePhoto   = 1
	TypeAudio   = 2
)

type Reply struct {
	ChatId    string
	ReplyType string
	Content   interface{}
}
```

### Building the replies

To easily construct the replies for the bot, use the `ReplyBuilder`.

```go
package example

import "github.com/detecc/detecctor-v2/model/reply"

func CreateReply() {
	builder := reply.NewReplyBuilder()
	replyMessage1 := builder.TypeMessage().ForChat("chatId1").WithContent("sampleContent").Build()
	replyMessage2 := builder.TypePhoto().ForChat("chatId2").WithContent(nil).Build()
}
```