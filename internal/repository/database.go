package repository

import (
	"database/sql"
	"fmt"
	"playground/internal/entities"
	"playground/internal/entities/database"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/driver/sqliteshim"
)

type Database struct {
	*bun.DB
}

func newDatabase(config database.Config) (*bun.DB, error) {
	switch config := config.(type) {
	case nil:
		const dsn = ":memory:"
		sqldb, err := sql.Open(sqliteshim.ShimName, dsn)
		if err != nil {
			return nil, err
		}
		return bun.NewDB(sqldb, sqlitedialect.New()), nil
	case database.PostgresConfig:
		dsn := config.DSN()
		sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
		db := bun.NewDB(sqldb, pgdialect.New())
		return db, nil
	case database.SQLiteConfig:
		if config.Filename == "" {
			return nil, fmt.Errorf("required field is empty: %s", "Filename")
		}
		sqldb, err := sql.Open(sqliteshim.ShimName, config.Filename)
		if err != nil {
			return nil, err
		}
		return bun.NewDB(sqldb, sqlitedialect.New()), nil
	}
	return nil, entities.ErrorUnsupportedConfig
}
