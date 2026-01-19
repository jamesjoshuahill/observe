package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/jamesjoshuahill/observe/internal/browser"
	"github.com/jamesjoshuahill/observe/internal/config"
	"github.com/jamesjoshuahill/observe/internal/tools"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "config":
			return runConfig()
		case "list":
			return runList()
		case "validate":
			return runValidate()
		}
	}

	return runOpen()
}

func runConfig() error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	configPath := config.Path()
	if err := ensureConfigDir(configPath); err != nil {
		return err
	}

	cmd := exec.Command(editor, configPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func ensureConfigDir(configPath string) error {
	dir := configPath[:strings.LastIndex(configPath, "/")]
	return os.MkdirAll(dir, 0755)
}

func runList() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	fmt.Println("Environments:")
	envNames := make([]string, 0, len(cfg.Environments))
	for name := range cfg.Environments {
		envNames = append(envNames, name)
	}
	sort.Strings(envNames)
	for _, name := range envNames {
		fmt.Printf("  %s\n", name)
	}

	fmt.Println("\nServices:")
	svcNames := make([]string, 0, len(cfg.Services))
	for name := range cfg.Services {
		svcNames = append(svcNames, name)
	}
	sort.Strings(svcNames)
	for _, name := range svcNames {
		envs := make([]string, 0, len(cfg.Services[name]))
		for env := range cfg.Services[name] {
			envs = append(envs, env)
		}
		sort.Strings(envs)
		fmt.Printf("  %s (%s)\n", name, strings.Join(envs, ", "))
	}

	return nil
}

func runValidate() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if err := cfg.Validate(); err != nil {
		return err
	}

	fmt.Println("Config is valid.")
	return nil
}

func runOpen() error {
	fs := flag.NewFlagSet("observe", flag.ExitOnError)
	service := fs.String("service", "", "Service name")
	env := fs.String("env", "", "Environment name")
	toolsFlag := fs.String("tools", "", "Comma-separated list of tools (default: all)")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: observe [command] [flags]\n\n")
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "  config    Open config file in $EDITOR\n")
		fmt.Fprintf(os.Stderr, "  list      List configured services and environments\n")
		fmt.Fprintf(os.Stderr, "  validate  Validate config file\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}

	if *service == "" || *env == "" {
		fs.Usage()
		return fmt.Errorf("--service and --env are required")
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	envConfig, err := cfg.GetEnvironment(*env)
	if err != nil {
		return err
	}

	svcConfig, err := cfg.GetServiceEnv(*service, *env)
	if err != nil {
		return err
	}

	var selectedTools []tools.Tool
	if *toolsFlag != "" {
		for _, name := range strings.Split(*toolsFlag, ",") {
			name = strings.TrimSpace(name)
			t, ok := tools.Get(name)
			if !ok {
				return fmt.Errorf("unknown tool: %s", name)
			}
			selectedTools = append(selectedTools, t)
		}
	} else {
		selectedTools = tools.All()
	}

	for _, t := range selectedTools {
		url, err := t.BuildURL(envConfig, svcConfig)
		if err != nil {
			var notConfigured tools.ErrNotConfigured
			if errors.As(err, &notConfigured) {
				fmt.Fprintf(os.Stderr, "warning: %s\n", err)
				continue
			}
			return err
		}

		fmt.Printf("Opening %s: %s\n", t.Name(), url)
		if err := browser.Open(url); err != nil {
			return err
		}
	}

	return nil
}
