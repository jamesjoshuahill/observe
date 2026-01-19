package pagerduty

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var incidentPathRegex = regexp.MustCompile(`^/incidents/([A-Z0-9]+)$`)

// ParseIncidentURL extracts the incident ID from a PagerDuty incident URL.
// Expected format: https://*.pagerduty.com/incidents/PXXXXXX
func ParseIncidentURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid PagerDuty URL: %w", err)
	}

	if u.Scheme != "https" {
		return "", fmt.Errorf("invalid PagerDuty URL: expected https scheme")
	}

	if !strings.HasSuffix(u.Host, ".pagerduty.com") && u.Host != "pagerduty.com" {
		return "", fmt.Errorf("invalid PagerDuty URL: expected *.pagerduty.com host")
	}

	matches := incidentPathRegex.FindStringSubmatch(u.Path)
	if matches == nil {
		return "", fmt.Errorf("invalid PagerDuty URL: expected format https://*.pagerduty.com/incidents/PXXXXXX")
	}

	return matches[1], nil
}
