package payload

import (
	"crypto/rand"
	"fmt"
	goCache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

const ErrGeneratingUuid = "unknown"

// Uuid creates unique identifier.
func Uuid() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Println(err)
		return "unknown"
	}

	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return uuid
}

// ParseIP separates the IP and Port of the address.
func ParseIP(addr string) (string, string) {
	ip, port, err := net.SplitHostPort(addr)
	if err != nil {
		log.Debugf("Address %s not valid", addr)
		return "", ""
	}

	return ip, port
}

// GeneratePayloadId Generates a UUID for an outbound Payload to map the response to the ChatId
func GeneratePayloadId(cache *goCache.Cache, payload *Payload, chatId string) {
	// Create a unique id for every server message
	log.Info("Generating new payloadId..")
	uuid := Uuid()
	log.Debugf("UUID: %s", uuid)
	if uuid == ErrGeneratingUuid {
		// Bad
		log.WithField("payload", payload).Errorf("uuid couldnt be generated")
		return
	}

	payload.Id = uuid
	// Set the payloadId to chatId mapping
	cache.Set(uuid, chatId, time.Minute*5)
}
