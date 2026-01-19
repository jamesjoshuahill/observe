package tools

import (
	"fmt"

	"github.com/jamesjoshuahill/observe/internal/config"
)

// Tool defines the interface for observability tools.
type Tool interface {
	Name() string
	BuildURL(envConfig *config.EnvironmentConfig, svcConfig *config.ServiceEnvConfig) (string, error)
}

var registry = map[string]Tool{
	"grafana":   &Grafana{},
	"kibana":    &Kibana{},
	"pagerduty": &PagerDuty{},
}

// Get returns a tool by name.
func Get(name string) (Tool, bool) {
	t, ok := registry[name]
	return t, ok
}

// All returns all registered tools.
func All() []Tool {
	return []Tool{
		registry["grafana"],
		registry["kibana"],
		registry["pagerduty"],
	}
}

// Names returns the names of all registered tools.
func Names() []string {
	return []string{"grafana", "kibana", "pagerduty"}
}

// ErrNotConfigured indicates a tool is not configured for a service/environment.
type ErrNotConfigured struct {
	Tool string
}

func (e ErrNotConfigured) Error() string {
	return fmt.Sprintf("%s not configured", e.Tool)
}
