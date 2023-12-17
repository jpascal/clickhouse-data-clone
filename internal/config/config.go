package config

type Database struct {
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
}

type Config struct {
	Source      Database `yaml:"source"`
	Destination Database `yaml:"destination"`
	Skip        []string `yaml:"skip"`
}
