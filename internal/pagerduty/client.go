package pagerduty

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

const apiBaseURL = "https://api.pagerduty.com"

// Incident represents the extracted metadata from a PagerDuty incident.
type Incident struct {
	Service     string
	Environment string
	RunbookURL  string
}

// Client is a PagerDuty API client.
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new PagerDuty API client.
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

// GetIncident fetches an incident by ID and extracts service/environment metadata.
func (c *Client) GetIncident(id string) (*Incident, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/incidents/%s", apiBaseURL, id), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Token token="+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching incident: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("PagerDuty API error (status %d): %s", resp.StatusCode, string(body))
	}

	var apiResp incidentResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return extractIncidentMetadata(&apiResp)
}

type incidentResponse struct {
	Incident struct {
		Body struct {
			Details map[string]interface{} `json:"details"`
		} `json:"body"`
		Description string `json:"description"`
		Alerts      []struct {
			Body struct {
				Details map[string]interface{} `json:"details"`
			} `json:"body"`
		} `json:"alerts"`
	} `json:"incident"`
}

func extractIncidentMetadata(resp *incidentResponse) (*Incident, error) {
	incident := &Incident{}

	// Try to extract from incident body details first
	details := resp.Incident.Body.Details

	// Also check alert details if incident details are empty
	if len(details) == 0 && len(resp.Incident.Alerts) > 0 {
		details = resp.Incident.Alerts[0].Body.Details
	}

	// Extract service
	if svc, ok := getStringField(details, "service"); ok {
		incident.Service = svc
	} else {
		return nil, fmt.Errorf("incident missing 'service' field in details")
	}

	// Extract environment
	if env, ok := getStringField(details, "environment"); ok {
		incident.Environment = env
	} else {
		return nil, fmt.Errorf("incident missing 'environment' field in details")
	}

	// Extract runbook URL (optional)
	if runbook, ok := getStringField(details, "runbook_url"); ok {
		incident.RunbookURL = runbook
	} else if runbook, ok := getStringField(details, "runbook"); ok {
		incident.RunbookURL = runbook
	} else {
		// Try to find URL in description
		incident.RunbookURL = extractURLFromText(resp.Incident.Description)
	}

	return incident, nil
}

func getStringField(m map[string]interface{}, key string) (string, bool) {
	if m == nil {
		return "", false
	}
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok && s != "" {
			return s, true
		}
	}
	return "", false
}

var urlRegex = regexp.MustCompile(`https?://[^\s<>"]+`)

func extractURLFromText(text string) string {
	matches := urlRegex.FindAllString(text, -1)
	for _, match := range matches {
		// Look for URLs that look like runbooks
		lower := strings.ToLower(match)
		if strings.Contains(lower, "runbook") ||
			strings.Contains(lower, "wiki") ||
			strings.Contains(lower, "confluence") ||
			strings.Contains(lower, "notion") ||
			strings.Contains(lower, "docs") {
			return match
		}
	}
	// Return first URL if no runbook-specific URL found
	if len(matches) > 0 {
		return matches[0]
	}
	return ""
}
