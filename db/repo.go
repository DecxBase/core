package db

import (
	"github.com/uptrace/bun"
)

type DataRepository[P any, M any] struct {
	db            *bun.DB
	tempDB        *bun.DB
	activeFilters []*QueryFilters
	schema        FilterSchemas
}

func NewRepository[P any, M any](db *bun.DB) *DataRepository[P, M] {
	return &DataRepository[P, M]{
		db:            db,
		schema:        make(FilterSchemas),
		activeFilters: make([]*QueryFilters, 0),
	}
}
