package repo

import (
	"fmt"
	"reflect"

	"github.com/DecxBase/core/types"
	"github.com/DecxBase/core/utils"
	"github.com/uptrace/bun"
)

type repoCacheData struct {
	Service       string
	Paginate      bool
	PerPage       int32
	KeywordFields []string
	ActionFilters map[types.ServiceActionType][]string
	KeywordKey    string
	PageKey       string
	PerPageKey    string
}

func ResolveCache(source any) (*repoCacheData, error) {
	typ := reflect.TypeOf(source).Elem()

	m := source.(types.ModelForService)
	if m == nil {
		return nil, fmt.Errorf("unable to resolve service model: %s", typ.String())
	}

	actions := make(map[types.ServiceActionType][]string)
	for _, ac := range []types.ServiceActionType{
		types.ServiceActionFindAll,
		types.ServiceActionFind,
		types.ServiceActionDelete,
	} {
		actions[ac] = m.ServiceActionToFilters(ac)
	}

	return &repoCacheData{
		Service:       m.ModelForService(),
		Paginate:      m.UsePagination(),
		PerPage:       m.PaginatePerPage(),
		KeywordFields: m.KeywordFields(),
		ActionFilters: actions,
		KeywordKey:    "keyword",
		PageKey:       "page",
		PerPageKey:    "per_page",
	}, nil
}

type DataFieldSchema struct {
	Expr      string
	Translate string
	Condition types.DataFilterCondition
	Multiples []string
	UseOR     bool
	GroupCond string
	Builder   func(bun.QueryBuilder, ...any) bun.QueryBuilder
}

func (s DataFieldSchema) Apply(q bun.QueryBuilder, key string, values ...any) bun.QueryBuilder {
	if s.Builder != nil {
		return s.Builder(q, values...)
	}

	if len(s.Expr) > 0 {
		return q.Where(s.Expr, values...)
	}

	if s.Multiples != nil {
		groupCond := " AND "
		if len(s.GroupCond) > 0 {
			groupCond = s.GroupCond
		}

		return q.WhereGroup(groupCond, func(qb bun.QueryBuilder) bun.QueryBuilder {
			for _, f := range s.Multiples {
				qb = utils.ApplyCondition(qb, s.Condition, s.UseOR, f, values...)
			}

			return qb
		})
	}

	theKey := key
	if len(s.Translate) > 0 {
		theKey = s.Translate
	}

	return utils.ApplyCondition(q, s.Condition, s.UseOR, theKey, values...)
}

func (r *DataRepository[P, M]) UseSchemaBuilder(cb func(types.ServiceActionType, bun.QueryBuilder) bun.QueryBuilder) *DataRepository[P, M] {
	r.schemaBuilder = cb
	return r
}

func (r *DataRepository[P, M]) UseSchemaDataReader(cb func(types.ServiceActionType, any) *types.RepoSchemaData) *DataRepository[P, M] {
	r.schemaDataReader = cb
	return r
}

func (r *DataRepository[P, M]) AddFields(names ...string) *DataRepository[P, M] {
	for _, name := range names {
		r.fieldSchemas[name] = &DataFieldSchema{}
	}

	return r
}

func (r *DataRepository[P, M]) AddFieldSchema(name string, schema ...*DataFieldSchema) *DataRepository[P, M] {
	if len(schema) > 0 {
		r.fieldSchemas[name] = schema[0]
	} else {
		r.fieldSchemas[name] = &DataFieldSchema{}
	}

	return r
}

func (r *DataRepository[P, M]) ResolveSchemaData(action types.ServiceActionType, source any) *types.RepoSchemaData {
	var data *types.RepoSchemaData

	if source != nil {
		if r.schemaDataReader != nil {
			data = r.schemaDataReader(action, source)
		} else {
			data = utils.ParseSchemaData(action, source)
		}
	}

	return data
}

func (r *DataRepository[P, M]) ResolveSchemaFilter(action types.ServiceActionType, data *types.RepoSchemaData) *DataFilters {
	filters := NewDataFilters()

	if data != nil && data.Filter != nil {
		fields := r.Cache.ActionFilters[action]

		kFields := r.Cache.KeywordFields
		kValue := data.Filter[r.Cache.KeywordKey]
		if kFields != nil && kValue != nil {
			filters = filters.Like(kValue, kFields...)
		}

		if fields != nil {
			filters = filters.WhereFunc(func(q bun.QueryBuilder) bun.QueryBuilder {
				for _, key := range fields {
					field := r.fieldSchemas[key]
					fValue := data.Filter[key]

					if field != nil && fValue != nil {
						q = field.Apply(q, key, fValue)
					}
				}

				return q
			})
		}
	}

	if r.schemaBuilder != nil {
		return filters.WhereFunc(func(qb bun.QueryBuilder) bun.QueryBuilder {
			return r.schemaBuilder(action, qb)
		})
	}

	return filters
}

func (r *DataRepository[P, M]) ResolveSchemaPaginate(query *bun.SelectQuery, data *types.RepoSchemaData) *bun.SelectQuery {
	if data != nil && r.Cache.Paginate {
		page, perPage := r.ExtractPaginationData(data)

		return query.Offset((int(page) - 1) * int(perPage)).Limit(int(perPage))
	}

	return query
}

func (r *DataRepository[P, M]) ExtractPaginationData(data *types.RepoSchemaData) (int, int) {
	if data != nil && r.Cache.Paginate {
		var page int = 1
		var perPage int = int(r.Cache.PerPage)

		if data.Paginate != nil {
			tPage := data.Paginate[r.Cache.PageKey]
			if tPage != nil {
				page = int(utils.StringToInt64(fmt.Sprintf("%v", tPage)))
			}

			tPerPage := data.Paginate[r.Cache.PerPageKey]
			if tPerPage != nil {
				perPage = int(utils.StringToInt64(fmt.Sprintf("%v", tPerPage)))
			}
		}

		if page < 1 {
			page = 1
		}

		return page, perPage
	}

	return 0, 0
}
