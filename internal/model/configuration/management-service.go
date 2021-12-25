package configuration

type (
	ManagementServiceConfiguration struct {
		MqttBroker    MqttBroker    `fig:"mqttBroker" validate:"required"`
		Database      Database      `fig:"database" validate:"required"`
		Observability Observability `fig:"observability"`
	}
)
