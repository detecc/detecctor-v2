package payload

import (
	"crypto/rand"
	"fmt"
	"github.com/detecc/detecctor-v2/internal/cache"
	"github.com/detecc/detecctor-v2/model/payload"
	"log"
	"net"
	"time"
)

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
		log.Println("Address not valid")
		return "", ""
	}
	return ip, port
}

// GeneratePayloadId Generates a UUID for an outbound Payload to map the response to the ChatId
func GeneratePayloadId(payload *payload.Payload, chatId string) {
	//create a unique id for every server message
	uuid := Uuid()
	log.Println("UUID:", uuid)
	if uuid == "" {
		// bad
		log.Println("uuid is empty")
		return
	}

	payload.Id = uuid
	//set the payload ID to chatId mapping
	cache.Memory().Set(uuid, chatId, time.Minute*5)
}
