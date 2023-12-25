package tasks

import (
	"context"
	"database/sql"
	"fmt"
	intConfig "github.com/jpascal/clickhouse-data-clone/internal/config"
	"github.com/jpascal/clickhouse-data-clone/internal/logging"
	"github.com/jpascal/clickhouse-data-clone/internal/modules"
	"github.com/pkg/errors"
	"slices"
	"time"
)

const tableSQL = `SELECT DISTINCT table_name, table_type FROM information_schema.tables where table_catalog = '%s' ORDER BY table_name ASC`
const insertSQL = `INSERT INTO FUNCTION remote('%s', '%s.%s', '%s', '%s') SELECT * FROM %s.%s`
const insertWithFilterSQL = `INSERT INTO FUNCTION remote('%s', '%s.%s', '%s', '%s') SELECT * FROM %s.%s WHERE %s`

const (
	BaseTable = `BASE TABLE`
	//ViewTable = TableType(`VIEW`)
)

type Table struct {
	Name string
	Type string
}

func Execute(ctx context.Context) error {
	startAt := time.Now()
	config := intConfig.GetCtx(ctx)
	logger := logging.GetCtx(ctx)

	logger.Info("connecting to source")
	source, err := modules.Database(config.Source)
	if err != nil {
		return errors.Wrap(err, "source")
	}

	logger.Info("connecting to destination")
	destination, err := modules.Database(config.Destination)
	if err != nil {
		return errors.Wrap(err, "destination")
	}

	tables, err := source.Query(ctx, fmt.Sprintf(tableSQL, config.Source.Name))
	if err != nil {
		return errors.Wrap(err, "source.tables")
	}

	logger.Info("start database coping")

	for tables.Next() {
		var table Table
		if err := tables.Scan(
			&table.Name, &table.Type,
		); err != nil {
			return errors.Wrap(err, "tables.scan")
		}

		tableLogger := logger.With("table", table.Name, "type", table.Type)

		if slices.Contains(config.Tables.Skip, table.Name) {
			tableLogger.Warn("skip table")
		}

		if len(config.Tables.Only) > 0 && !slices.Contains(config.Tables.Only, table.Name) {
			tableLogger.Warn("skip table")
			continue
		}

		if config.Recreate {
			tableLogger.Info("re-create table")
			row := source.QueryRow(ctx, fmt.Sprintf(`SHOW CREATE TABLE %s.%s`, config.Source.Name, table.Name))
			var ddl string
			if err = row.Scan(&ddl); err != nil {
				return errors.Wrap(err, "source.ddl")
			}
			err = destination.Exec(ctx, fmt.Sprintf("DROP TABLE IF EXISTS %s.%s", config.Destination.Name, table.Name))
			if !errors.Is(err, sql.ErrNoRows) && err != nil {
				return errors.Wrap(err, "destination.drop")
			}
			err = destination.Exec(ctx, ddl)
			if !errors.Is(err, sql.ErrNoRows) && err != nil {
				return errors.Wrap(err, "destination.create")
			}
		}
		if config.Truncate && table.Type == BaseTable {
			tableLogger.Info("truncate table")
			err = destination.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s.%s", config.Destination.Name, table.Name))
			if !errors.Is(err, sql.ErrNoRows) && err != nil {
				return errors.Wrap(err, "destination.truncate")
			}
		}

		if table.Type == BaseTable {
			startCopyAt := time.Now()
			tableLogger.Info("start table coping")
			if filter, ok := config.Tables.Filters[table.Name]; ok {
				tableLogger.Info("start table coping data by filter")
				err = source.Exec(ctx, fmt.Sprintf(insertWithFilterSQL,
					config.Destination.Host,
					config.Destination.Name,
					table.Name,
					config.Destination.User,
					config.Destination.Password,
					config.Source.Name,
					table.Name,
					filter,
				))
			} else {
				tableLogger.Info("start table coping all data")
				err = source.Exec(ctx, fmt.Sprintf(insertSQL,
					config.Destination.Host,
					config.Destination.Name,
					table.Name,
					config.Destination.User,
					config.Destination.Password,
					config.Source.Name,
					table.Name,
				))
			}
			if !errors.Is(err, sql.ErrNoRows) && err != nil {
				if config.Tables.Force {
					tableLogger.Error(errors.Wrap(err, "destination.insert").Error())
					continue
				}
				return errors.Wrap(err, "destination.insert")
			}
			tableLogger.With("elapsed", time.Since(startCopyAt)).Info("table copied successful")
		}
	}
	logger.With("elapsed", time.Since(startAt)).Info("database copied successful")
	return nil
}
