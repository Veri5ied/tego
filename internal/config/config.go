package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Hooks map[string]interface{} `json:"hooks" toml:"hooks" yaml:"hooks"`
}

var configFiles = []string{
	".tegorc.json",
	".tegorc",
	"tego.json",
	"tego.toml",
	"tego.yaml",
	"tego.yml",
}

func ConfigExists() bool {
	for _, filename := range configFiles {
		if fileExists(filename) {
			return true
		}
	}
	return false
}

func GetConfigPath() string {
	for _, filename := range configFiles {
		if fileExists(filename) {
			return filename
		}
	}
	return ""
}

func LoadConfig() (*Config, error) {
	for _, filename := range configFiles {
		if !fileExists(filename) {
			continue
		}

		data, err := os.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("error reading %s: %v", filename, err)
		}

		config := &Config{}

		// Determine format and parse
		switch {
		case isJSON(filename):
			err = json.Unmarshal(data, config)
		case isTOML(filename):
			err = toml.Unmarshal(data, config)
		case isYAML(filename):
			err = yaml.Unmarshal(data, config)
		default:
			// JSON as default
			err = json.Unmarshal(data, config)
		}

		if err != nil {
			return nil, fmt.Errorf("error parsing %s: %v", filename, err)
		}

		return config, nil
	}

	return nil, fmt.Errorf("no configuration file found. Run 'tego init' to create one")
}

func (c *Config) GetCommands(hookName string) ([]string, error) {
	hookCmd, exists := c.Hooks[hookName]
	if !exists {
		return nil, fmt.Errorf("hook '%s' not configured", hookName)
	}

	var commands []string
	switch v := hookCmd.(type) {
	case string:
		commands = []string{v}
	case []interface{}:
		for _, cmd := range v {
			if cmdStr, ok := cmd.(string); ok {
				commands = append(commands, cmdStr)
			}
		}
	default:
		return nil, fmt.Errorf("invalid hook configuration for '%s'", hookName)
	}

	return commands, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func isJSON(filename string) bool {
	return strings.HasSuffix(filename, ".json") || filename == ".tegorc"
}

func isTOML(filename string) bool {
	return strings.HasSuffix(filename, ".toml")
}

func isYAML(filename string) bool {
	return strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml")
}
