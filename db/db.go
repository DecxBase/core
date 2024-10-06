package db

import (
	"context"
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

// CreateSchema creates database schema
//
//	[]interface{}{
//		(*model.Category)(nil),
//	}
func CreateSchema(ctx context.Context, db *bun.DB, models []interface{}) error {
	for _, model := range models {
		_, err := db.NewCreateTable().Model(model).Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// DropSchema removes database schema
func DropSchema(ctx context.Context, db *bun.DB, models []interface{}) error {
	for _, model := range models {
		_, err := db.NewDropTable().Model(model).Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
