# Project Context

CLI tool that opens observability tools (Grafana, Kibana, PagerDuty) in the browser for a given service and environment.

## Build & Test

```bash
go build ./cmd/observe    # Build
go test ./...             # Run tests (none yet)
go mod tidy               # Tidy dependencies
```

## Project Structure

```text
cmd/observe/main.go       # CLI entry point, subcommands, flag parsing
internal/
  config/config.go        # Config loading from ~/.config/observe/config.yaml
  browser/browser.go      # Cross-platform URL opening (macOS/Linux)
  tools/
    tools.go              # Tool interface and registry
    grafana.go            # Grafana URL builder
    kibana.go             # Kibana URL builder
    pagerduty.go          # PagerDuty URL builder
```

## Architecture

- **Config**: YAML file defines environments (base URLs) and services (tool-specific IDs/queries per env)
- **Tools**: Each tool implements `Tool` interface with `Name()` and `BuildURL()` methods
- **Browser**: Uses `open` (macOS) or `xdg-open` (Linux) to open URLs

## Adding a New Tool

1. Create `internal/tools/<toolname>.go` implementing the `Tool` interface
2. Add to registry in `internal/tools/tools.go`
3. Add config fields to `EnvironmentConfig` and `ServiceEnvConfig` in `internal/config/config.go`

## Dependencies

- `gopkg.in/yaml.v3` - YAML parsing
- No other external dependencies
