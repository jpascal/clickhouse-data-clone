package main

import (
	"context"
	"github.com/jpascal/clickhouse-data-clone/internal/config"
	"github.com/jpascal/clickhouse-data-clone/internal/logging"
	"github.com/jpascal/clickhouse-data-clone/internal/modules"
	"github.com/jpascal/clickhouse-data-clone/internal/tasks"
	"github.com/jpascal/clickhouse-data-clone/internal/variables"
	"github.com/urfave/cli/v2"
	"log/slog"
	"os"
)

func main() {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx = logging.SetCtx(ctx, logger)

	cfg, err := modules.Config()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	ctx = config.SetCtx(ctx, cfg)

	app := cli.App{
		Name:     variables.Banner("Data Clone"),
		Version:  variables.Version,
		Usage:    "clone data from one to another ClickHouse databases",
		HideHelp: false,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:      "config",
				Value:     "confix.yaml",
				EnvVars:   []string{"CONFIG"},
				TakesFile: false,
			},
		},
		Commands: []*cli.Command{
			{
				Name: "run",
				Action: func(cliCtx *cli.Context) error {
					return tasks.Execute(cliCtx.Context)
				},
			},
		},
	}
	if err := app.RunContext(ctx, os.Args); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
