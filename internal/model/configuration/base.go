package configuration

type (
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

	Database struct {
		Database string `fig:"database" default:"detecctor"`
		Host     string `fig:"host" default:"localhost"`
		Username string `fig:"username" validate:"required"`
		Password string `fig:"password" validate:"required"`
		Port     int    `fig:"port" default:"27017"`
	}
)
