package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the top-level configuration.
type Config struct {
	PagerDutyAPIKey string                       `yaml:"pagerduty_api_key"`
	Environments    map[string]EnvironmentConfig `yaml:"environments"`
	Services        map[string]ServiceConfig     `yaml:"services"`
}

// EnvironmentConfig holds base URLs for each tool in an environment.
type EnvironmentConfig struct {
	Grafana string `yaml:"grafana"`
	Kibana  string `yaml:"kibana"`
}

// ServiceConfig maps environment names to tool-specific settings.
type ServiceConfig map[string]ServiceEnvConfig

// ServiceEnvConfig holds tool-specific IDs and queries for a service in an environment.
type ServiceEnvConfig struct {
	GrafanaDashboardID string `yaml:"grafana_dashboard_id"`
	KibanaQuery        string `yaml:"kibana_query"`
}

// Path returns the default config file path.
func Path() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "observe", "config.yaml")
}

// Load reads and parses the config file.
func Load() (*Config, error) {
	return LoadFrom(Path())
}

// LoadFrom reads and parses a config file from the given path.
func LoadFrom(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	return &cfg, nil
}

// Validate checks that the config is well-formed.
func (c *Config) Validate() error {
	if len(c.Environments) == 0 {
		return fmt.Errorf("no environments configured")
	}
	if len(c.Services) == 0 {
		return fmt.Errorf("no services configured")
	}

	for svcName, svcConfig := range c.Services {
		for envName := range svcConfig {
			if _, ok := c.Environments[envName]; !ok {
				return fmt.Errorf("service %q references unknown environment %q", svcName, envName)
			}
		}
	}

	return nil
}

// GetServiceEnv returns the service configuration for a given service and environment.
func (c *Config) GetServiceEnv(service, env string) (*ServiceEnvConfig, error) {
	svc, ok := c.Services[service]
	if !ok {
		return nil, fmt.Errorf("unknown service: %s", service)
	}

	envConfig, ok := svc[env]
	if !ok {
		return nil, fmt.Errorf("service %q has no configuration for environment %q", service, env)
	}

	return &envConfig, nil
}

// GetEnvironment returns the environment configuration for a given environment.
func (c *Config) GetEnvironment(env string) (*EnvironmentConfig, error) {
	envConfig, ok := c.Environments[env]
	if !ok {
		return nil, fmt.Errorf("unknown environment: %s", env)
	}
	return &envConfig, nil
}
