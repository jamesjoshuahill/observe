package tools

import (
	"fmt"
	"strings"

	"github.com/jamesjoshuahill/observe/internal/config"
)

// PagerDuty implements the Tool interface for PagerDuty.
type PagerDuty struct{}

func (p *PagerDuty) Name() string {
	return "pagerduty"
}

func (p *PagerDuty) BuildURL(envConfig *config.EnvironmentConfig, svcConfig *config.ServiceEnvConfig) (string, error) {
	if envConfig.PagerDuty == "" {
		return "", ErrNotConfigured{Tool: "pagerduty"}
	}
	if svcConfig.PagerDutyServiceID == "" {
		return "", ErrNotConfigured{Tool: "pagerduty"}
	}

	baseURL := strings.TrimSuffix(envConfig.PagerDuty, "/")
	return fmt.Sprintf("%s/services/%s", baseURL, svcConfig.PagerDutyServiceID), nil
}
