# observe

A CLI that opens observability tools (Grafana, Kibana, PagerDuty) for a given service and environment.

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

# List configured services and environments
observe list

# Validate config file
observe validate

# Open config in $EDITOR
observe config
```

## Configuration

Config file location: `~/.config/observe/config.yaml`

### Example config

```yaml
environments:
  prod:
    grafana: "https://grafana.example.com"
    kibana: "https://kibana.example.com"
    pagerduty: "https://example.pagerduty.com"
  staging:
    grafana: "https://grafana-staging.example.com"
    kibana: "https://kibana-staging.example.com"

services:
  api:
    prod:
      grafana_dashboard_id: "abc123"
      kibana_query: "service:api"
      pagerduty_service_id: "PXXXXXX"
    staging:
      grafana_dashboard_id: "abc123-staging"
      kibana_query: "service:api AND env:staging"
  web:
    prod:
      grafana_dashboard_id: "xyz789"
      kibana_query: "service:web"
```

### URL construction

| Tool      | URL pattern                                                    |
|-----------|----------------------------------------------------------------|
| Grafana   | `{base}/d/{dashboard_id}`                                      |
| Kibana    | `{base}/app/discover#/?_g=()&_a=(query:(query_string:(query:'{query}')))` |
| PagerDuty | `{base}/services/{service_id}`                                 |

## Supported platforms

- macOS (uses `open`)
- Linux (uses `xdg-open`)
