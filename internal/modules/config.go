package modules

import "github.com/jpascal/clickhouse-data-clone/internal/config"

func Config() (*config.Config, error) {
	return Register("config", func(s string) (*config.Config, error) {
		return config.Load("config.yaml")
	})
}

//func Tasks() (*config.Tasks, error) {
//	return Register("tasks", func(s string) (*config.Tasks, error) {
//		cfg, err := Config()
//		if err != nil {
//			return nil, err
//		}
//		return &cfg.Tasks, err
//	})
//}
