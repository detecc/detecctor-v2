package configuration

type (
	BaseServiceConfiguration struct {
		MqttBroker MqttBroker `fig:"mqttBroker" validate:"required"`
		Mongo      Database   `fig:"mongo" validate:"required"`
	}

	NotificationServiceConfiguration struct {
		BaseServiceConfiguration
		Bot BotConfiguration `fig:"bot" validate:"required"`
	}

	PluginServiceConfiguration struct {
		BaseServiceConfiguration
		PluginConfiguration PluginConfiguration `fig:"plugins" validate:"required"`
	}

	ManagementServiceConfiguration struct {
		BaseServiceConfiguration
	}

	MqttBroker struct {
		Host     string `fig:"host" default:"localhost"`
		Port     int    `fig:"port" default:"1883"`
		Username string `fig:"username" validate:"required"`
		Password string `fig:"password" validate:"required"`
		ClientId string `fig:"clientId" validate:"required"`
		Tls      TLS    `fig:"tls"`
	}

	TLS struct {
		IsEnabled  bool   `fig:"isEnabled"`
		CACertPath string `fig:"CACertPath"`
		KeyPath    string `fig:"keyPath"`
	}

	PluginConfiguration struct {
		PluginDir string   `fig:"dir" validate:"required,dir"`
		Plugins   []string `fig:"list"`
	}

	BotConfiguration struct {
		Type  string `fig:"type" validate:"required"`
		ID    string `fig:"id" validate:"required"`
		Token string `fig:"token" validate:"required"`
	}

	Database struct {
		Database string `fig:"database" default:"detecctor"`
		Host     string `fig:"host" default:"localhost"`
		Username string `fig:"username" validate:"required"`
		Password string `fig:"password" validate:"required"`
		Port     int    `fig:"port" default:"27017"`
	}
)
