package management

import (
	. "github.com/detecc/detecctor-v2/internal/mqtt"
	"github.com/eclipse/paho.mqtt.golang"
)

const (
	ChatAuth                 = Topic("chat/+/auth")
	ChatSetLang              = Topic("chat/+/lang/set")
	ClientRegister           = Topic("client/+/register")
	ClientHeartbeat          = Topic("client/+/heartbeat")
	ClientRegisterReplyTopic = Topic("client/+/register/reply")
)

var ChatAuthHandler = func(client mqtt.Client, message mqtt.Message) {

}

var SetLanguageHandler = func(client mqtt.Client, message mqtt.Message) {

}

var ClientRegisterHandler = func(client mqtt.Client, message mqtt.Message) {

}

var PluginExecuteRegisterHandler = func(client mqtt.Client, message mqtt.Message) {

}

var NotificationServiceRegisterHandler = func(client mqtt.Client, message mqtt.Message) {

}
