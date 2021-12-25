package config

import (
	configuration2 "github.com/detecc/detecctor-v2/internal/model/configuration"
	"github.com/detecc/detecctor-v2/pkg/cache"
	"github.com/kkyr/fig"
	goCache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"path/filepath"
)

const (
	PluginServiceConfiguration       = "PluginServiceConfiguration"
	NotificationServiceConfiguration = "NotificationServiceConfiguration"
	ManagementServiceConfiguration   = "ManagementServiceConfiguration"

	dockerConfigDir = "/detecctor-v2/configuration/"
	confDir         = "../../config"
)

func GetServiceConfiguration(config interface{}, cache *goCache.Cache, filePath, cacheKey string) {
	var (
		err error
	)

	// Load the config
	err = fig.Load(config,
		fig.File(filepath.Base(filePath)),
		fig.Dirs(filepath.Dir(filePath), dockerConfigDir, confDir, "."),
	)
	if err != nil {
		log.WithError(err).Fatal("Unable to load configuration")
	}

	if cache != nil {
		// Cache the config
		cache.Set(cacheKey, &config, goCache.NoExpiration)
	}
}

// GetNotificationServiceConfiguration get the configuration from the configuration file and store the configuration in the cache
func GetNotificationServiceConfiguration(filePath string) *configuration2.NotificationServiceConfiguration {
	var (
		config              configuration2.NotificationServiceConfiguration
		memory              = cache.Memory()
		isFound             bool
		cachedConfiguration interface{}
	)

	// Check if the configuration is cached
	cachedConfiguration, isFound = memory.Get(NotificationServiceConfiguration)
	if isFound {
		return cachedConfiguration.(*configuration2.NotificationServiceConfiguration)
	}

	// Get configuration
	GetServiceConfiguration(&config, memory, filePath, NotificationServiceConfiguration)

	return &config
}

// GetPluginServiceConfiguration get the configuration from the configuration file and store the configuration in the cache
func GetPluginServiceConfiguration(filePath string) *configuration2.PluginServiceConfiguration {
	var (
		config              configuration2.PluginServiceConfiguration
		memory              = cache.Memory()
		isFound             bool
		cachedConfiguration interface{}
	)

	// Check if the configuration is cached
	cachedConfiguration, isFound = memory.Get(PluginServiceConfiguration)
	if isFound {
		return cachedConfiguration.(*configuration2.PluginServiceConfiguration)
	}

	// Get configuration
	GetServiceConfiguration(&config, memory, filePath, PluginServiceConfiguration)

	return &config
}

// GetManagementServiceConfiguration get the configuration from the configuration file and store the configuration in the cache
func GetManagementServiceConfiguration(filePath string) *configuration2.PluginServiceConfiguration {
	var (
		config              configuration2.PluginServiceConfiguration
		memory              = cache.Memory()
		isFound             bool
		cachedConfiguration interface{}
	)

	// Check if the configuration is cached
	cachedConfiguration, isFound = memory.Get(ManagementServiceConfiguration)
	if isFound {
		return cachedConfiguration.(*configuration2.PluginServiceConfiguration)
	}

	// Get configuration
	GetServiceConfiguration(&config, memory, filePath, ManagementServiceConfiguration)

	return &config
}
