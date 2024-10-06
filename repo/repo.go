package repo

import (
	"github.com/DecxBase/core/types"
	"github.com/uptrace/bun"
)

type DataRepository[P any, M types.ModelForService] struct {
	db    *bun.DB
	Cache *repoCacheData

	schemaDataReader func(types.ServiceActionType, any) *types.RepoSchemaData
	schemaBuilder    func(types.ServiceActionType, bun.QueryBuilder) bun.QueryBuilder
	fieldSchemas     map[string]*DataFieldSchema

	tempDB        *bun.DB
	keepFilters   bool
	activeFilters []*DataFilters
}

func NewRepository[P any, M types.ModelForService](db *bun.DB) *DataRepository[P, M] {
	mm := new(M)
	cache, _ := ResolveCache(mm)

	return &DataRepository[P, M]{
		db:           db,
		Cache:        cache,
		fieldSchemas: make(map[string]*DataFieldSchema),

		keepFilters:   false,
		activeFilters: make([]*DataFilters, 0),
	}
}
