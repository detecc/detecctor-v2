package configuration

type (
	NotificationServiceConfiguration struct {
		MqttBroker    MqttBroker       `fig:"mqttBroker" validate:"required"`
		Database      Database         `fig:"database" validate:"required"`
		Bot           BotConfiguration `fig:"bot" validate:"required"`
		Observability Observability    `fig:"observability"`
	}

	BotConfiguration struct {
		Type  string `fig:"type" validate:"required"`
		ID    string `fig:"id" validate:"required"`
		Token string `fig:"token" validate:"required"`
	}
)
