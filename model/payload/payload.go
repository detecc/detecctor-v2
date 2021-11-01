package payload

type (
	// Payload is used to transfer data and determine the status, target client and plugin.
	Payload struct {
		// Id is used to uniquely identify the request to the response.
		Id string `json:"Id" validate:"required"`
		// ServiceNodeKey is used to determine the target client for the data.
		ServiceNodeKey string      `json:"serviceNodeKey" validate:"required"`
		Data           interface{} `json:"data,omitempty"`
		Command        string      `json:"command" validate:"required"`
		Success        bool        `json:"success"`
		Error          string      `json:"error,omitempty"`
	}

	Option func(*Payload)
)

// SetError indicate that something went wrong with the payload
func (payload *Payload) SetError(err error) {
	if err != nil {
		payload.Success = false
		payload.Error = err.Error()
	}
}

// NewPayload creates a new payload with the options provided
func NewPayload(opts ...Option) Payload {
	var payload = &Payload{
		Id:             "",
		ServiceNodeKey: "",
		Data:           nil,
		Command:        "",
		Success:        false,
		Error:          "",
	}

	for _, opt := range opts {
		opt(payload)
	}

	return *payload
}

// ForClient add a target client for the payload
func ForClient(serviceNodeKey string) Option {
	return func(payload *Payload) {
		payload.ServiceNodeKey = serviceNodeKey
	}
}

// ForCommand add a command which the payload should execute or is originating from
func ForCommand(command string) Option {
	return func(payload *Payload) {
		payload.Command = command
	}
}

// WithData sets the payload data
func WithData(data interface{}) Option {
	return func(payload *Payload) {
		payload.Data = data
	}
}

// Successful set that the payload is successful
func Successful() Option {
	return func(payload *Payload) {
		payload.Error = ""
		payload.Success = true
	}
}

// WithError sets an error for the payload
func WithError(err error) Option {
	return func(payload *Payload) {
		if err == nil {
			return
		}
		payload.SetError(err)
	}
}
