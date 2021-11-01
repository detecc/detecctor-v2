# Authentication

## User authentication

Before you can issue commands to Detecctor, you must you have access to the server resources by authenticating yourself
with a token.

The authentication process is simple:

1. Issue `/auth` command to the bot.
2. The server will generate a unique token for the chat.
3. The token will be put in the **logs/console/file** on the host.
4. The user will have 5 minutes to claim the token. If not claimed, it will be deleted and the user must be
   re-authenticated.
5. Issue the `/auth` command with the authentication key as an argument. Example: `/auth xaA3fVhg5qwefgg2g6b4`

## Client authentication

Before the client is authorized to communicate with the server, receive and execute commands, it must be authenticated
with the system. By default, the client's first message after successfully establishing the connection with the server
is sending an authentication request.

In both the server and client's configuration files contain an `authPassword` attribute. The client will be authorized
if the `authPassword` in both configurations match.