# observe

A CLI that opens observability tools (Grafana, Kibana) for a given service and environment.

Based on a tool I co-developed with colleagues at Form3 that helped us reduce the fatigue of monitoring our services across 10+ environments.

## Installation

```bash
go install github.com/jamesjoshuahill/observe/cmd/observe@latest
```

## Usage

```bash
# Open all configured tools for a service/environment
observe --service api --env prod

# Open specific tools only
observe --service api --env prod --tools grafana,kibana

# Incident response: open tools based on PagerDuty alert
observe --alert https://example.pagerduty.com/incidents/P1234567

# List configured services and environments
observe list

# Validate config file
observe validate

# Open config in $EDITOR
observe config
```

### Incident response mode

The `--alert` flag enables incident response workflow:

1. Fetches incident details from PagerDuty API
2. Extracts `service` and `environment` from alert metadata
3. Opens all configured observability tools for that service/environment
4. Opens runbook URL if present in the alert details

This requires `pagerduty_api_key` in your config and alerts with `service` and `environment` fields in the incident body details.

## Configuration

Config file location: `~/.config/observe/config.yaml`

### Example config

```yaml
pagerduty_api_key: "u+XXXXX..."  # Required for --alert flag

environments:
  prod:
    grafana: "https://grafana.example.com"
    kibana: "https://kibana.example.com"
  staging:
    grafana: "https://grafana-staging.example.com"
    kibana: "https://kibana-staging.example.com"

services:
  api:
    prod:
      grafana_dashboard_id: "abc123"
      kibana_query: "service:api"
    staging:
      grafana_dashboard_id: "abc123-staging"
      kibana_query: "service:api AND env:staging"
  web:
    prod:
      grafana_dashboard_id: "xyz789"
      kibana_query: "service:web"
```

### URL construction

| Tool    | URL pattern                                                    |
|---------|----------------------------------------------------------------|
| Grafana | `{base}/d/{dashboard_id}`                                      |
| Kibana  | `{base}/app/discover#/?_g=()&_a=(query:(query_string:(query:'{query}')))` |

## Supported platforms

- macOS (uses `open`)
- Linux (uses `xdg-open`)
