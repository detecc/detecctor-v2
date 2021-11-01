# Detecctor-v2 core concepts

Detecctor-v2 is an improved version of the previous version Detecctor; instead of using TCP connection, the
server-client communication now uses MQTT for improved connectivity and performance.

## Architecture

Another core difference in the v2 version is in the server architecture. Instead of using a plain TCP server, we're now
using the microservice architecture with 3 core services:

- The notification service
- The plugin service
- The management service

The services also use MQTT as a communication layer, to simplify and improve responsiveness. It also makes more sense
from the scalability and performance perspective - if you've got a lot of plugins and/or clients, you simply just scale
the plugin microservice.

### Plugin service

The plugin microservice does all the processing. It creates payloads and sends them to clients and if/when a client
responds, it processes the data sent back from the client. It then forwards the response to the Notification service if
needed.

### Notification service

Notification service is essentially a bridge between various bots and other notification services and the end user. It
logs and processes the message, creates a command and sends it to the Plugin service. When the Plugin service sends a
response, it is sent back to the user.

### Management service

The management service manages user and client access to different plugins and commands. It interacts with the EMQ X api
to manage the access to MQTT topics. Clients should only be able to communicate with the Plugin and Management services
and not interact with the notification service.