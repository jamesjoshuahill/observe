package tools

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/jamesjoshuahill/observe/internal/config"
)

// Kibana implements the Tool interface for Kibana.
type Kibana struct{}

func (k *Kibana) Name() string {
	return "kibana"
}

func (k *Kibana) BuildURL(envConfig *config.EnvironmentConfig, svcConfig *config.ServiceEnvConfig) (string, error) {
	if envConfig.Kibana == "" {
		return "", ErrNotConfigured{Tool: "kibana"}
	}
	if svcConfig.KibanaQuery == "" {
		return "", ErrNotConfigured{Tool: "kibana"}
	}

	baseURL := strings.TrimSuffix(envConfig.Kibana, "/")
	encodedQuery := url.QueryEscape(svcConfig.KibanaQuery)
	return fmt.Sprintf("%s/app/discover#/?_g=()&_a=(query:(query_string:(query:'%s')))", baseURL, encodedQuery), nil
}
