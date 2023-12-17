package config

import (
	"github.com/goccy/go-yaml"
	"io"
	"os"
)

func Load(path string) (*Config, error) {
	configFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var configData string
	if data, err := io.ReadAll(configFile); err != nil {
		return nil, err
	} else {
		configData = string(data)
	}
	configData = os.ExpandEnv(configData)
	cfg := Config{}
	return &cfg, yaml.Unmarshal([]byte(configData), &cfg)
}
