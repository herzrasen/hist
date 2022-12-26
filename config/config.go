package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type Config struct {
	Excludes []string `yaml:"excludes"`
}

func Load(path string) (*Config, error) {
	path, err := resolvePath(path)
	if err != nil {
		return nil, fmt.Errorf("config:Config:Load:resolvePath: %w", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config:Config:load: read file: %w", err)
	}
	var config Config
	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("config:Config:load: unmarshal: %w", err)
	}
	return &config, nil
}

func resolvePath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		userDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("config:Config:resolvePath: %w", err)
		}
		return strings.Join([]string{userDir, strings.TrimPrefix(path, "~/")}, "/"), nil
	}
	return path, nil
}

func (c *Config) IsExcluded(s string) bool {
	if c.Excludes == nil {
		return false
	}
	for _, exclude := range c.Excludes {
		if s == exclude {
			return true
		}
	}
	return false
}
