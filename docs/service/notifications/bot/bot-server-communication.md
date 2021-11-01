# The bot and bot proxy communication

The server-proxy-bot communication occurs on channels through dedicated listeners. The message that comes from a chat
service's Bot is sent through the `chan ProxyMessage` channel to the bot `Proxy`. The proxy processes the message and
constructs a `Command` struct that is then sent to the server through the `chan Command` channel.

The server processes the command, executes the plugins and makes a `Reply`. The reply is sent through a `chan Reply`
channel and is passed to the proxy, which passes it to the bot.

## The Bot

The Bot interface enables the project to support multiple bots with two-way communication between the user and the
server. You can change the type of bot you want to have in
the [settings](../../../service/notifications/configuration.md#notification-service-configuration).

Check the [bot guide](adding-bots.md) for more info about the bots.

## Bot Proxy

The Bot Proxy manages the communication between the Bot and the server. It also logs and persists all the necessary data
to the database. A Bot, which is created for a specific type, is passed to the Bot Proxy. The proxy then automatically
starts the bot and its listeners.
