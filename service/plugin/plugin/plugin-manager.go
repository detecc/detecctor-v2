package plugin

import (
	"fmt"
	"github.com/detecc/detecctor-v2/internal/model/configuration"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"plugin"
	"sync"
)

var pluginManager Manager

func init() {
	once := sync.Once{}
	once.Do(func() {
		GetPluginManager()
	})
}

type (
	// ManagerImpl is a manager for the plugins. It stores and maps the plugins to the command.
	ManagerImpl struct {
		plugins sync.Map
	}

	Manager interface {
		HasPlugin(name string) bool
		AddPlugin(name string, plugin Handler)
		GetPlugin(name string) (Handler, error)
		LoadPlugins(pluginConfig configuration.PluginConfiguration)
	}
)

// Register a cmd to the manager.
func Register(name string, plugin Handler) {
	GetPluginManager().AddPlugin(name, plugin)
}

// GetPluginManager gets the cmd manager instance (singleton).
func GetPluginManager() Manager {
	if pluginManager == nil {
		pluginManager = &ManagerImpl{plugins: sync.Map{}}
	}
	return pluginManager
}

// HasPlugin Check if the cmd exists in the manager.
func (pm *ManagerImpl) HasPlugin(name string) bool {
	_, exists := pm.plugins.Load(name)
	return exists
}

// AddPlugin Add a cmd to the manager.
func (pm *ManagerImpl) AddPlugin(name string, plugin Handler) {
	log.WithField("name", name).Debug("Adding cmd to manager")
	if !pm.HasPlugin(name) {
		pm.plugins.Store(name, plugin)
	}
}

// GetPlugin returns the cmd, if found.
func (pm *ManagerImpl) GetPlugin(name string) (Handler, error) {
	mPlugin, exists := pm.plugins.Load(name)
	if exists {
		return mPlugin.(Handler), nil
	}

	return nil, fmt.Errorf("cmd doesnt exist")
}

// LoadPlugins Load the plugins from the folder, specified in the configuration file.
func (pm *ManagerImpl) LoadPlugins(pluginConfig configuration.PluginConfiguration) {
	log.Info("Loading plugins..")

	for _, pluginFromList := range pluginConfig.Plugins {

		err := filepath.Walk(pluginConfig.PluginDir, func(path string, info os.FileInfo, err error) error {
			// if the name matches, try to load the cmd
			if !info.IsDir() && info.Name() == pluginFromList+".so" {
				log.Infof("Loading cmd: %s", pluginFromList)
				_, err = plugin.Open(path)
				return err
			}

			return nil
		})
		if err != nil {
			log.WithError(err).Errorf("Cannot load cmd %s from list", pluginFromList)
		}

	}
}
