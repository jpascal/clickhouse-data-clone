package config

import "time"

type Database struct {
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Timeouts struct {
		Read time.Duration `json:"read"`
		Dial time.Duration `json:"dial"`
	} `yaml:"timeouts"`
}

type TableFilter string

type Tables struct {
	Force   bool                   `yaml:"force"`
	Only    []string               `yaml:"only"`
	Skip    []string               `yaml:"skip"`
	Filters map[string]TableFilter `yaml:"filters"`
}

type Config struct {
	Source      Database `yaml:"source"`
	Destination Database `yaml:"destination"`
	Truncate    bool     `yaml:"truncate"`
	Recreate    bool     `yaml:"recreate"`
	Tables      Tables   `yaml:"tables"`
}
