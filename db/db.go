package db

import (
	"context"
	"database/sql"
	"fmt"
	"runtime"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func GetDB(connector *pgdriver.Connector) *bun.DB {
	maxOpenConns := 4 * runtime.GOMAXPROCS(0)

	sqldb := sql.OpenDB(connector)
	sqldb.SetMaxOpenConns(maxOpenConns)
	sqldb.SetMaxIdleConns(maxOpenConns)

	err := sqldb.Ping()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %s", err))
	}

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
