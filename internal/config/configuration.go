package config

import (
	"flag"
	"fmt"
	cache2 "github.com/detecc/detecctor-v2/internal/cache"
	"github.com/detecc/detecctor-v2/model/configuration"
	"github.com/kkyr/fig"
	"github.com/patrickmn/go-cache"
	"log"
	"os"
	"path/filepath"
)

const (
	PluginService       = Service("plugin")
	NotificationService = Service("notification")
	ManagementService   = Service("management")
)

type Service string

// GetFlags get the program flags and store them in the cache
func GetFlags(service Service) {
	var (
		memory              = cache2.Memory()
		configFileName      *string
		defaultFileName     = ""
		cacheKey            = ""
		workingDirectory, _ = os.Getwd()
		configurationPath   = fmt.Sprintf("%s", workingDirectory)
	)

	// Get the paths from arguments
	configurationFileFormatFlag := flag.String("config-format", "yaml", "Format of the configuration files (YAML, JSON or TOML)")

	switch service {
	case PluginService:
		defaultFileName = "plugin-config"
		cacheKey = "pluginServiceConfigFilePath"
		break
	case NotificationService:
		defaultFileName = "notification-config"
		cacheKey = "notificationServiceConfigFilePath"
		break
	case ManagementService:
		defaultFileName = "management-config"
		cacheKey = "managementServiceConfigFilePath"
		break
	default:
		return
	}

	fileName := fmt.Sprintf("%s/configs/%s.%s", configurationPath, defaultFileName, *configurationFileFormatFlag)
	configFileName = flag.String("config-file", fileName, "Path of the configuration file")
	flag.Parse()

	memory.Set(cacheKey, *configFileName, cache.NoExpiration)
}

// GetNotificationServiceConfiguration get the configuration from the configuration file and store the configuration in the cache
func GetNotificationServiceConfiguration() *configuration.NotificationServiceConfiguration {
	var (
		config                configuration.NotificationServiceConfiguration
		err                   error
		configurationFilePath string
		memory                = cache2.Memory()
		isFound               bool
		cachedConfiguration   interface{}
	)

	cachedConfiguration, isFound = memory.Get("notificationServiceConfiguration")
	if isFound {
		return cachedConfiguration.(*configuration.NotificationServiceConfiguration)
	}

	configurationPath, isFound := memory.Get("notificationServiceConfigFilePath")
	if isFound {
		configurationFilePath = configurationPath.(string)
	} else {
		log.Fatal("No configuration file path found!")
	}

	err = fig.Load(&config,
		fig.File(filepath.Base(configurationFilePath)),
		fig.Dirs(filepath.Dir(configurationFilePath)),
	)
	if err != nil {
		log.Fatal(err)
	}

	memory.Set("notificationServiceConfiguration", &config, cache.NoExpiration)

	return &config
}

// GetPluginServiceConfiguration get the configuration from the configuration file and store the configuration in the cache
func GetPluginServiceConfiguration() *configuration.PluginServiceConfiguration {
	var (
		config                configuration.PluginServiceConfiguration
		err                   error
		configurationFilePath string
		memory                = cache2.Memory()
		isFound               bool
		cachedConfiguration   interface{}
	)

	cachedConfiguration, isFound = memory.Get("pluginServiceConfiguration")
	if isFound {
		return cachedConfiguration.(*configuration.PluginServiceConfiguration)
	}

	configurationPath, isFound := memory.Get("pluginServiceConfigFilePath")
	if isFound {
		configurationFilePath = configurationPath.(string)
	} else {
		log.Fatal("No configuration file path found!")
	}

	err = fig.Load(&config,
		fig.File(filepath.Base(configurationFilePath)),
		fig.Dirs(filepath.Dir(configurationFilePath)),
	)
	if err != nil {
		log.Fatal(err)
	}

	memory.Set("pluginServiceConfiguration", &config, cache.NoExpiration)

	return &config
}

// GetManagementServiceConfiguration get the configuration from the configuration file and store the configuration in the cache
func GetManagementServiceConfiguration() *configuration.PluginServiceConfiguration {
	var (
		config                configuration.PluginServiceConfiguration
		err                   error
		configurationFilePath string
		memory                = cache2.Memory()
		isFound               bool
		cachedConfiguration   interface{}
	)

	cachedConfiguration, isFound = memory.Get("managementServiceConfiguration")
	if isFound {
		return cachedConfiguration.(*configuration.PluginServiceConfiguration)
	}

	configurationPath, isFound := memory.Get("managementServiceConfigFilePath")
	if isFound {
		configurationFilePath = configurationPath.(string)
	} else {
		log.Fatal("No configuration file path found!")
	}

	err = fig.Load(&config,
		fig.File(filepath.Base(configurationFilePath)),
		fig.Dirs(filepath.Dir(configurationFilePath)),
	)
	if err != nil {
		log.Fatal(err)
	}

	memory.Set("managementServiceConfiguration", &config, cache.NoExpiration)

	return &config
}
