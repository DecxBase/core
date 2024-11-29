package repo

import (
	"fmt"
	"strconv"

	"github.com/DecxBase/core/types"
	"github.com/DecxBase/core/utils"
	"github.com/uptrace/bun"
)

type dataFilterCallback = func(bun.QueryBuilder) bun.QueryBuilder
type selectFilterCallback = func(*bun.SelectQuery) *bun.SelectQuery

type DataFilterCondition struct {
	Expr           string
	Compare        types.DataFilterCondition
	Callback       dataFilterCallback
	SelectCallback selectFilterCallback
	Fields         []string
	Value          any
	Values         []any
}

type DataFilters struct {
	page       int
	perPage    int
	conditions []DataFilterCondition
}

func (f *DataFilters) add(cond DataFilterCondition) *DataFilters {
	f.conditions = append(
		f.conditions,
		cond,
	)

	return f
}

func (f *DataFilters) Clear() *DataFilters {
	f.conditions = make([]DataFilterCondition, 0)
	return f
}

func (f *DataFilters) SetPage(page int) *DataFilters {
	f.page = page
	return f
}

func (f *DataFilters) SetPerPage(perPage int) *DataFilters {
	f.perPage = perPage
	return f
}

func (f *DataFilters) Merge(filters ...*DataFilters) *DataFilters {
	for _, fl := range filters {
		f.conditions = append(
			f.conditions,
			fl.conditions...,
		)
	}

	return f
}

func (f DataFilters) AddExpr(expr string, comp types.DataFilterCondition, values ...any) *DataFilters {
	return f.add(DataFilterCondition{
		Expr:    expr,
		Compare: comp,
		Values:  values,
	})
}

func (f DataFilters) Where(expr string, values ...any) *DataFilters {
	return f.AddExpr(expr, types.DataFilterWhere, values...)
}

func (f DataFilters) Like(value any, fields ...string) *DataFilters {
	return f.add(DataFilterCondition{
		Compare: types.DataFilterWhereLike,
		Fields:  fields,
		Value:   value,
	})
}

func (f DataFilters) WhereFunc(cb dataFilterCallback) *DataFilters {
	return f.add(DataFilterCondition{
		Compare:  types.DataFilterFunc,
		Callback: cb,
	})
}

func (f DataFilters) Relation(expr string, cb selectFilterCallback) *DataFilters {
	return f.add(DataFilterCondition{
		Compare:        types.DataFilterRelation,
		SelectCallback: cb,
	})
}

func (f DataFilters) WherePK() *DataFilters {
	return f.AddExpr("pk", types.DataFilterWherePK)
}

func (f DataFilters) ApplyBuilder(query bun.QueryBuilder) bun.QueryBuilder {
	for _, cond := range f.conditions {
		switch cond.Compare {
		case types.DataFilterWhere:
			query = query.Where(cond.Expr, cond.Values...)
		case types.DataFilterWherePK:
			query = query.WherePK()
		case types.DataFilterWhereLike:
			query = query.WhereGroup(" AND ", func(qb bun.QueryBuilder) bun.QueryBuilder {
				for _, f := range cond.Fields {
					qb = utils.ApplyCondition(qb, types.DataFilterWhereLike, true, f, cond.Value)
				}

				return qb
			})
		case types.DataFilterFunc:
			query = cond.Callback(query)
		}
	}

	return query
}

func (f DataFilters) ApplySelectBuilder(query *bun.SelectQuery) *bun.SelectQuery {
	for _, cond := range f.conditions {
		switch cond.Compare {
		case types.DataFilterRelation:
			query = query.Relation(cond.Expr, cond.SelectCallback)
		}
	}

	return query
}

func (f DataFilters) ApplyPaginate(query *bun.SelectQuery) *bun.SelectQuery {
	return query.Limit(f.perPage).Offset((f.page - 1) * f.perPage)
}

func (f *DataFilters) Load(data map[string]any) *DataFilters {
	f.Clear()

	for key, value := range data {
		if key == "page" {
			page, _ := strconv.Atoi(value.(string))
			f.SetPage(page)
		} else if key == "per_page" {
			perPage, _ := strconv.Atoi(value.(string))
			f.SetPerPage(perPage)
		} else {
			f.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}

	return f
}

func NewDataFilters() *DataFilters {
	qf := &DataFilters{
		page:       1,
		perPage:    10,
		conditions: make([]DataFilterCondition, 0),
	}

	return qf
}

func MapToDataFilters(data map[string]any) *DataFilters {
	return NewDataFilters().Load(data)
}
