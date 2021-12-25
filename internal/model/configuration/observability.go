package configuration

type (
	Observability struct {
		OTel    OTel    `fig:"tracing" required:""`
		Metrics Metrics `fig:"metrics" required:""`
		Logging Logging `fig:"logging" required:""`
	}

	OTel struct {
		Address  string `fig:"address" required:""`
		AuthType string `fig:"auth" default:""`
		TLS      TLS    `fig:"tls"`
	}

	Metrics struct {
		Address  string `fig:"address" required:""`
		Endpoint string `fig:"url" default:"/metrics"`
		AuthType string `fig:"auth" default:""`
		Username string `fig:"username" default:""`
		Password string `fig:"password" default:""`
		TLS      TLS    `fig:"tls"`
	}

	Logging struct {
		Type   []string `fig:"type" validate:"required"` // file, remote, console
		Format string   `fig:"format" default:"syslog"`  // syslog, json, etc
		Host   string   `fig:"host"`
		Port   int      `fig:"port" default:"1514"`
	}
)
