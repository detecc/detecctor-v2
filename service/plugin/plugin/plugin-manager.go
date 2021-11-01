package plugin

import (
	"fmt"
	"github.com/detecc/detecctor-v2/model/configuration"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"plugin"
	"sync"
)

var pluginManager *Manager

func init() {
	once := sync.Once{}
	once.Do(func() {
		GetPluginManager()
	})
}

// Manager is a manager for the plugins. It stores and maps the plugins to the command.
type Manager struct {
	plugins sync.Map
}

// Register a plugin to the manager.
func Register(name string, plugin Handler) {
	GetPluginManager().AddPlugin(name, plugin)
}

// GetPluginManager gets the plugin manager instance (singleton).
func GetPluginManager() *Manager {
	if pluginManager == nil {
		pluginManager = &Manager{plugins: sync.Map{}}
	}
	return pluginManager
}

// HasPlugin Check if the plugin exists in the manager.
func (pm *Manager) HasPlugin(name string) bool {
	_, exists := pm.plugins.Load(name)
	return exists
}

// AddPlugin Add a plugin to the manager.
func (pm *Manager) AddPlugin(name string, plugin Handler) {
	log.WithField("name", name).Debug("Adding plugin to manager")
	if !pm.HasPlugin(name) {
		pm.plugins.Store(name, plugin)
	}
}

// GetPlugin returns the plugin, if found.
func (pm *Manager) GetPlugin(name string) (Handler, error) {
	mPlugin, exists := pm.plugins.Load(name)
	if exists {
		return mPlugin.(Handler), nil
	}

	return nil, fmt.Errorf("plugin doesnt exist")
}

// LoadPlugins Load the plugins from the folder, specified in the configuration file.
func (pm *Manager) LoadPlugins(pluginConfig configuration.PluginConfiguration) {
	log.Info("Loading plugins..")

	for _, pluginFromList := range pluginConfig.Plugins {

		err := filepath.Walk(pluginConfig.PluginDir, func(path string, info os.FileInfo, err error) error {
			// if the name matches, try to load the plugin
			if !info.IsDir() && info.Name() == pluginFromList+".so" {
				log.Infof("Loading plugin: %s", pluginFromList)
				_, err = plugin.Open(path)
				return err
			}

			return nil
		})

		if err != nil {
			log.Errorf("error loading plugin from list: %v", err)
		}
	}
}
