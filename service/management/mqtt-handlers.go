package management

import (
	"fmt"
	"github.com/detecc/detecctor-v2/database"
	"github.com/detecc/detecctor-v2/internal/cache"
	. "github.com/detecc/detecctor-v2/internal/mqtt"
	payload2 "github.com/detecc/detecctor-v2/model/payload"
	"github.com/detecc/detecctor-v2/service/management/auth"
	log "github.com/sirupsen/logrus"
)

const (
	ChatAuth                 = Topic("chat/+/auth")
	ChatSetLang              = Topic("chat/+/lang/set")
	ClientRegister           = Topic("client/+/register")
	ClientHeartbeat          = Topic("client/+/heartbeat")
	ClientRegisterReplyTopic = Topic("client/+/register/reply")
)

var ChatAuthHandler = func(client Client, topicIds []string, payloadId uint16, payload interface{}, err error) {
	chatId := topicIds[0]
	token := ""

	log.Debugf("Authorizing a chat:%s", chatId)

	if database.GetChatRepository().IsChatAuthorized(nil, chatId) {
		client.Publish("", nil)
		return
	}

	// Check if the token is in the cache and if it matches the provided token
	cachedTokenId := fmt.Sprintf("auth-token-%s", chatId)
	cachedToken, isFound := cache.Memory().Get(cachedTokenId)
	if isFound && cachedToken.(string) == token {
		err := database.GetChatRepository().AuthorizeChat(nil, chatId)
		if err != nil {
			log.WithFields(log.Fields{
				"chatId": chatId,
				"token":  token,
			}).Errorf("Error authorizing chat: %v", err)

			client.Publish("", nil)
			return
		}

		cache.Memory().Delete(cachedTokenId)
	}

	if !isFound && token == "" {
		// Generate a token
		auth.GenerateChatAuthenticationToken(chatId)
		client.Publish("", nil)
	}

}

var SetLanguageHandler = func(client Client, topicIds []string, payloadId uint16, payload interface{}, err error) {
	chatId := topicIds[0]

	err = database.GetChatRepository().SetLanguage(chatId, "lang")
	if err != nil {
		log.WithFields(log.Fields{
			"chatId":   chatId,
			"language": "lang",
		}).Errorf("Error updating the language: %v", err)

		client.Publish("", fmt.Sprintf("An error occured while setting the language: %v.", err))
		return
	}

	client.Publish("", fmt.Sprintf("Successfully set the language to: %s.", err))
}

var ClientRegisterHandler = func(client Client, topicIds []string, payloadId uint16, payload interface{}, err error) {
	clientId := topicIds[0]
	clientPayload := payload.(payload2.Payload)

	if clientPayload.Data == nil {
		log.Warnf("Payload data is empty; client %s cannot be authorized", clientId)
		return
	}

	if clientPayload.Data.(string) != "" {

	}

	// Try to authorize the client
	err = database.GetClientRepository().AuthorizeClient(nil, clientId, clientPayload.ServiceNodeKey)
	if err != nil {
		log.WithFields(log.Fields{
			"payload":  clientPayload,
			"clientId": clientId,
		}).Errorf("Error updating the client authorization status: %v", err)
		return
	}
}

var PluginExecuteRegisterHandler = func(client Client, topicIds []string, payloadId uint16, payload interface{}, err error) {
	//pluginName := topicIds[0]
	//todo
	// Forward to plugin service; check access rights and the client status

	if database.GetChatRepository().IsChatAuthorized(nil, "") &&
		database.GetClientRepository().IsClientAuthorized(nil, "") {

	}

}

var NotificationServiceRegisterHandler = func(client Client, topicIds []string, payloadId uint16, payload interface{}, err error) {

	//Todo
}
