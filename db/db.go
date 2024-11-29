package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"runtime"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/schema"
)

func PrepareSqlDB(db *sql.DB) *sql.DB {
	maxOpenConns := 4 * runtime.GOMAXPROCS(0)

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxOpenConns)

	err := db.Ping()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %s", err))
	}

	return db
}

func GetSqlDB(connector driver.Connector) *sql.DB {
	return PrepareSqlDB(OpenDB(connector))
}

func Transform(sqldb *sql.DB, dialect schema.Dialect) *bun.DB {
	db := bun.NewDB(
		sqldb, dialect,
		bun.WithDiscardUnknownColumns(),
	)
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	return db
}

func TransformConnector(connector driver.Connector, dialect schema.Dialect) *bun.DB {
	return Transform(GetSqlDB(connector), dialect)
}

func Open(driver string, dsn string) (*sql.DB, error) {
	return sql.Open(driver, dsn)
}

func OpenDB(connector driver.Connector) *sql.DB {
	return sql.OpenDB(connector)
}
