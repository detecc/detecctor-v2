package metrics

const (
	metricPrefix = "ocpp_service"

	/*--------------- API labels ---------------*/

	APIRequests  = metricPrefix + "api_requests"
	APIResponses = metricPrefix + "api_responses"
	APIErrors    = metricPrefix + "api_errors"

	/*--------------- OCPP labels ---------------*/

	OCPPRequests  = metricPrefix + "ocpp_requests"
	OCPPResponses = metricPrefix + "ocpp_responses"
	OCPPErrors    = metricPrefix + "ocpp_errors"

	/*--------------- Websocket labels ---------------*/

	ConnectedClients = metricPrefix + "currently_connected_clients"
)

type Registry struct {
}
