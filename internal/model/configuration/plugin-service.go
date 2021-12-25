package configuration

type (
	PluginServiceConfiguration struct {
		MqttBroker          MqttBroker          `fig:"mqttBroker" validate:"required"`
		Database            Database            `fig:"database" validate:"required"`
		PluginConfiguration PluginConfiguration `fig:"plugins" validate:"required"`
		Observability       Observability       `fig:"observability"`
	}

	PluginConfiguration struct {
		PluginDir string   `fig:"dir" validate:"required,dir"`
		Plugins   []string `fig:"list"`
	}
)
