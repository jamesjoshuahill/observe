package tools

import (
	"fmt"
	"strings"

	"github.com/jamesjoshuahill/observe/internal/config"
)

// Grafana implements the Tool interface for Grafana dashboards.
type Grafana struct{}

func (g *Grafana) Name() string {
	return "grafana"
}

func (g *Grafana) BuildURL(envConfig *config.EnvironmentConfig, svcConfig *config.ServiceEnvConfig) (string, error) {
	if envConfig.Grafana == "" {
		return "", ErrNotConfigured{Tool: "grafana"}
	}
	if svcConfig.GrafanaDashboardID == "" {
		return "", ErrNotConfigured{Tool: "grafana"}
	}

	baseURL := strings.TrimSuffix(envConfig.Grafana, "/")
	return fmt.Sprintf("%s/d/%s", baseURL, svcConfig.GrafanaDashboardID), nil
}
