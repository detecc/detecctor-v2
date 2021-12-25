package cache

import (
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

var mem *cache.Cache

func init() {
	once := sync.Once{}
	once.Do(func() {
		Memory()
	})
}
func NewCache() *cache.Cache {
	log.Info("Creating cache..")
	return cache.New(10*time.Minute, 5*time.Minute)
}

func Memory() *cache.Cache {
	if mem == nil {
		mem = NewCache()
	}

	return mem
}
