package db

import (
	"fmt"
	"strconv"

	"github.com/uptrace/bun"
)

type queryFilterCallback = func(bun.QueryBuilder) bun.QueryBuilder

type QueryFilterCondition struct {
	Expr     string
	Compare  string
	Callback queryFilterCallback
	Values   []any
}

type QueryFilters struct {
	page       int
	perPage    int
	conditions []QueryFilterCondition
}

func (f *QueryFilters) add(cond QueryFilterCondition) *QueryFilters {
	f.conditions = append(
		f.conditions,
		cond,
	)

	return f
}

func (f *QueryFilters) Clear() *QueryFilters {
	f.conditions = make([]QueryFilterCondition, 0)
	return f
}

func (f *QueryFilters) SetPage(page int) *QueryFilters {
	f.page = page
	return f
}

func (f *QueryFilters) SetPerPage(perPage int) *QueryFilters {
	f.perPage = perPage
	return f
}

func (f QueryFilters) AddExpr(expr string, comp string, values ...any) *QueryFilters {
	return f.add(QueryFilterCondition{
		Expr:    expr,
		Compare: comp,
		Values:  values,
	})
}

func (f QueryFilters) Where(expr string, values ...any) *QueryFilters {
	return f.AddExpr(expr, "where", values...)
}

func (f QueryFilters) WhereFunc(cb queryFilterCallback) *QueryFilters {
	return f.add(QueryFilterCondition{
		Compare:  "func",
		Callback: cb,
	})
}

func (f QueryFilters) WherePK() *QueryFilters {
	return f.AddExpr("pk", "where_pk")
}

func (f QueryFilters) ApplyBuilder(query bun.QueryBuilder) bun.QueryBuilder {
	for _, cond := range f.conditions {
		switch cond.Compare {
		case "where":
			query = query.Where(cond.Expr, cond.Values...)
		case "where_pk":
			query = query.WherePK()
		case "func":
			query = cond.Callback(query)
		}
	}

	return query
}

func (f QueryFilters) ApplyPaginate(query *bun.SelectQuery) *bun.SelectQuery {
	return query.Limit(f.perPage).Offset((f.page - 1) * f.perPage)
}

func (f *QueryFilters) Load(data map[string]any) *QueryFilters {
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

func NewQueryFilters() *QueryFilters {
	qf := &QueryFilters{
		page:       1,
		perPage:    10,
		conditions: make([]QueryFilterCondition, 0),
	}

	return qf
}
