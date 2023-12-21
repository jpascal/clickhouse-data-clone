package config

type Database struct {
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
}

type Tables struct {
	Only []string `yaml:"only"`
	Skip []string `yaml:"skip"`
}

type Config struct {
	Source      Database `yaml:"source"`
	Destination Database `yaml:"destination"`
	Truncate    bool     `yaml:"truncate"`
	Recreate    bool     `yaml:"recreate"`
	Tables      Tables   `yaml:"tables"`
}
