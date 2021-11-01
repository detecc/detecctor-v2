# Detecctor-v2

Detecctor-v2 is a ⚡fast, customizable 🖥️ client management and monitoring platform. It uses various 🤖 chatbots as a 📲
notification service. It is designed for use with 🔌 plugins, which enable total control over the functionality of both
server and the client(s). All you do is issue a command and let the 🔌 plugin deal with the rest. You can include provided
plugins, write your own or include plugins from the community.

## 🔧 Configuration

Before running the services, check out the [plugin service configuration](/docs/service/notifications/configuration.md)
and [notification service configuration](/docs/service/plugin/configuration.md) guides.

## 🤖 Supported bots

| Chat service | Supported |
| :----: | :----: |
| Telegram | ✔ |
| Slack | Planned |
| Discord | Planned |

## 🏃 Running the services

### Using 🐳 Docker or docker-compose

The provided _docker-compose_ file will run the MongoDB and the services.

```bash
docker-compose up -d
``` 

### Standalone

```bash
go build main.go 
./main # use --help to get the flags 
```

### 📢 Note:

The project is still under development. More features coming soon.