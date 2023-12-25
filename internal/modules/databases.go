package modules

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/jpascal/clickhouse-data-clone/internal/config"
	"github.com/pkg/errors"
	"strings"
	"time"
)

func Database(databaseConfig config.Database) (driver.Conn, error) {
	return Register[driver.Conn](strings.Join([]string{databaseConfig.Name, databaseConfig.User, databaseConfig.Host}, "/"), func(s string) (driver.Conn, error) {
		conn, err := clickhouse.Open(&clickhouse.Options{
			Addr:            []string{databaseConfig.Host},
			ReadTimeout:     time.Hour * 24,
			ConnMaxLifetime: time.Hour * 24,
			Auth: clickhouse.Auth{
				Database: databaseConfig.Name,
				Username: databaseConfig.User,
				Password: databaseConfig.Password,
			},

			ClientInfo: clickhouse.ClientInfo{
				Products: []struct {
					Name    string
					Version string
				}{
					{Name: "data-clone", Version: "0.0.1"},
				},
			},
			Debugf: func(format string, v ...interface{}) {
				fmt.Printf(format, v)
			},
		})
		if err != nil {
			return nil, errors.Wrap(err, "open")
		}
		if err := conn.Ping(context.Background()); err != nil {
			return nil, errors.Wrap(err, "ping")
		}
		return conn, nil
	})
}
