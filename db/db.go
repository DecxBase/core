package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"runtime"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func GetSqlDB(connector driver.Connector) *sql.DB {
	maxOpenConns := 4 * runtime.GOMAXPROCS(0)

	db := sql.OpenDB(connector)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxOpenConns)

	err := db.Ping()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %s", err))
	}

	return db
}

func GetBunDB(connector driver.Connector) *bun.DB {
	sqldb := GetSqlDB(connector)

	db := bun.NewDB(
		sqldb, pgdialect.New(),
		bun.WithDiscardUnknownColumns(),
	)
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	return db
}

func NewPGConnector(dsn string) *pgdriver.Connector {
	return pgdriver.NewConnector(pgdriver.WithDSN(dsn))
}
