package auth

import (
	"crypto/rand"
	"fmt"
	"github.com/detecc/detecctor-v2/internal/cache"
	log "github.com/sirupsen/logrus"
	"time"
)

// GenerateChatAuthenticationToken Generate an authorization token for a chat and log it. The token is cached and expires after 5 minutes.
func GenerateChatAuthenticationToken(chatId string) {
	log.WithField("chatId", chatId).Debug("Generating a new authentication token")

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.WithField("chatId", chatId).Error("Error generating a new authentication token")
		return
	}
	authToken := fmt.Sprintf("%x", b)
	cache.Memory().Set(fmt.Sprintf("auth-token-%s", chatId), authToken, time.Minute*5)
	log.Info(authToken)
}
