package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	intConfig "github.com/jpascal/clickhouse-data-clone/internal/config"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"
)

func main() {
	config, err := intConfig.Load("config.yaml")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	source, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{config.Source.Host},
		Auth: clickhouse.Auth{
			Database: config.Source.Name,
			Username: config.Source.User,
			Password: config.Source.Password,
		},
		ClientInfo: clickhouse.ClientInfo{
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "mover", Version: "0.1"},
			},
		},
		Debugf: func(format string, v ...interface{}) {
			fmt.Printf(format, v)
		},
	})
	if err != nil {
		panic(err)
	}

	destination, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{config.Destination.Host},
		Auth: clickhouse.Auth{
			Database: config.Destination.Name,
			Username: config.Destination.User,
			Password: config.Destination.Password,
		},
		ClientInfo: clickhouse.ClientInfo{
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "mover", Version: "0.1"},
			},
		},
		Debugf: func(format string, v ...interface{}) {
			fmt.Printf(format, v)
		},
	})
	if err != nil {
		panic(err)
	}

	if err := source.Ping(ctx); err != nil {
		var exception *clickhouse.Exception
		if errors.As(err, &exception) {
			fmt.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		panic(err)
	}

	rows, err := source.Query(ctx, fmt.Sprintf("SHOW tables FROM %s", config.Source.Name))
	if err != nil {
		log.Fatal(err)
	}

	startAt := time.Now()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	for rows.Next() {
		var (
			table          string
			migrateTableAt = time.Now()
		)
		if err := rows.Scan(
			&table,
		); err != nil {
			log.Fatal(err)
		}
		tableLogger := logger.With("table", table)
		//if config.Destination.Recreate {
		row := source.QueryRow(ctx, fmt.Sprintf(`SHOW CREATE TABLE %s.%s`, config.Source.Name, table))
		var ddl string
		if err = row.Scan(&ddl); err != nil {
			tableLogger.Error(err.Error())
			os.Exit(1)
		}

		err = destination.Exec(ctx, fmt.Sprintf("DROP TABLE IF EXISTS %s.%s", config.Destination.Name, table))
		if !errors.Is(err, sql.ErrNoRows) && err != nil {
			tableLogger.Error(err.Error())
			os.Exit(1)
		}

		err = destination.Exec(ctx, ddl)
		if !errors.Is(err, sql.ErrNoRows) && err != nil {
			tableLogger.Info(ddl)
			tableLogger.Error(err.Error())
			os.Exit(1)
		}

		if strings.Index(strings.ToUpper(ddl), "MATERIALIZED") == 0 {
			tableLogger.Info("skip materialized")
			continue
		}

		err = source.Exec(ctx, fmt.Sprintf(`INSERT INTO FUNCTION remote('%s', '%s.%s', '%s', '%s') SELECT * FROM %s.%s;`,
			config.Destination.Host,
			config.Destination.Name,
			table,
			config.Destination.User,
			config.Destination.Password,
			config.Source.Name,
			table,
		))
		if !errors.Is(err, sql.ErrNoRows) && err != nil {
			tableLogger.Error(err.Error())
			os.Exit(1)
		}
		tableLogger.With("elapsed", time.Since(migrateTableAt)).Info("migrated")
	}
	logger.With("elapsed", time.Since(startAt)).Info("migrated")
}
