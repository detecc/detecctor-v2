package auth

import (
	"crypto/rand"
	"fmt"
	goCache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"time"
)

// GenerateChatAuthenticationToken Generate an authorization token for a chat and log it. The token is cached and expires after 5 minutes.
func GenerateChatAuthenticationToken(cache *goCache.Cache, chatId string) {
	log.WithField("chatId", chatId).Debug("Generating a new authentication token")

	token := GenerateToken()
	if token == "" {
		return
	}

	cache.Set(fmt.Sprintf("auth-token-%s", chatId), token, time.Minute*5)
	log.Infof("Generated chat authentication token: %s", token)
}

func GenerateToken() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.WithError(err).Errorf("Error generating a new token")
		return ""
	}

	return fmt.Sprintf("%x", b)
}
