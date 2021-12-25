# Plugins

Plugins are used to customize the behaviour of the server and client to the specific needs of the user. There are a few
bundled plugins, which provide core functionality in the [Detecc-core](https://github.com/detecc/detecc-core)
repository. The plugin should have a `plugins.Handler` interface methods implemented in order to achieve desired
functionality.

All the plugins should be located in a `dir`, specified in
the [configuration file](../../../config/plugin-config.yaml).

The Bot shall receive and forward a command, which will find the corresponding plugin and invoke the `Execute` method of
the plugin. The `Execute` method returns an error and an array of `Payload`s that will be sent to the clients. After
sending the message and receiving the response from the client, the `Response` method will be invoked with the client's
response (`Payload.Data`) and will create a `Reply` struct to send back to the Bot Chat.

The `plugin.Handler` is shown below ([source file](../../../services/plugin/plugin/plugins.go)):

```go
package plugin

import (
	"github.com/detecc/detecctor-v2/model/reply"
	"github.com/detecc/detecctor-v2/model/payload"
)

type (
	Handler interface {

		// Response is called when the clients have responded and should
		// return a string to send as a reply to the bot
		Response(payload payload.Payload) reply.Reply

		// Execute method is called when the bot command matches GetCmdName's result.
		// The bot passes the string arguments to the method.
		// The execute method must return Payload array ready to be sent to the clients.
		Execute(args ...string) ([]payload.Payload, error)

		// GetMetadata returns the metadata about the cmd.
		GetMetadata() Metadata
	}

	// Metadata is used to determine the role of a cmd registered in the PluginManager.
	Metadata struct {

		// The Type of the cmd will determine the behaviour of the server and execution of the cmd(s).
		Type string

		// The Middleware list is used to determine, if the cmd has any middleware to execute.
		// Will be skipped if the cmd itself is registered as middleware.
		Middleware []string
	}
)
```

## Registering the plugins

You can register your plugins using:

```go
package main

import (
	"github.com/detecc/detecctor-v2/model/plugin"
)

func init() {
	example := &YourHandlerImplementation{}
	plugin.Register("/cmd-command", example)
	//or
	plugin.GetPluginManager().AddPlugin("/cmd-command", example)
}
```

## Translating replies

If you want to, you can translate the message you want to sent back to the chat. Firstly, add a new message in init
function by calling:

```go
package example

import (
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/detecc/detecctor-v2/internal/i18n"
)

func AddNewDefaultMessage() {
	// example of adding a new default message
	i18n.AddDefaultMessage(goi18n.Message{
		ID:    "Hello",
		Other: "Hello, World!",
	})
}
```

Then you need to extract messages using the go-i18n commands. Check out [go-i18n](https://github.com/nicksnyder/go-i18n)
for more insight into messages and translations.

To actually translate the message returned from the plugin, there are two options:

1. Using `i18n.Localize(lang string, messageId string, data map[string]interface{}, plural interface{})`
2. Using `server.TranslateReplyMessage(chatId int64, content interface{})`, which fetches the default language from the
   database and uses the `Localize` function by casting the content.

By using the recommended method `server.TranslateReplyMessage`, the `content` should be a map, created with an API:

```go
package example

import (
	"github.com/detecc/detecctor-v2/model/reply"
	"github.com/detecc/detecctor-v2/internal/i18n"
	"github.com/detecc/detecctor/server"
	"log"
)

func TranslateAndSendMessage() {
	// creates a new translation map with options
	translationMap := i18n.NewTranslationMap("messageId", i18n.AddData("key", "value"), i18n.WithPlural(1))

	// translate the desired message for the chat 
	message, err := server.TranslateReplyMessage("chatId", translationMap)
	if err != nil {
		return
	}
	log.Println(message)

	// send the message to the bot
	server.SendMessageToChat("chatId", reply.TypeMessage, message)
}
```

## Documenting plugins

Documentation of the plugins is important for further development and to better understand the logic of the plugin.

Your plugin documentation should contain:

1. What the plugin does and a brief introduction to the logic
2. The command and its arguments
3. Example(s) of the command call
4. Configuration file(s), if any are necessary
    - with a brief explanation of the attributes
    - default values, if any apply
5. The structure of the `Payload.Data`, if the plugin communicates with the client

## Plugin example

```go
package main

import (
	"log"
	"github.com/detecc/detecctor-v2/model/plugin"
	"github.com/detecc/detecctor-v2/model/reply"
	. "github.com/detecc/detecctor-v2/model/payload"
)

func init() {
	example := &Example{}
	plugin.Register("/example", example)
}

type Example struct {
	plugin.Handler
}

func (e Example) Response(payload Payload) reply.Reply {
	log.Println(payload)
	builder := reply.NewReplyBuilder()
	return builder.TypeMessage().WithContent("test").ForChat("chatId").Build()
}

func (e Example) Execute(args ...string) ([]Payload, error) {
	log.Println(args)
	return []Payload{}, nil
}

func (e Example) GetMetadata() plugin.Metadata {
	return plugin.Metadata{
		Type:       plugin.PluginTypeServerClient,
		Middleware: []string{},
	}
}
```