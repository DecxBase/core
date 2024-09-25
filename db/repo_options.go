package db

import (
	"fmt"

	"github.com/uptrace/bun"
)

type RepoCrudOptions struct {
	Data    map[string]any
	Columns []string
	Assign  []string
	Expr    string
}

func (o RepoCrudOptions) GetExpr() string {
	if len(o.Expr) > 0 {
		return o.Expr
	}

	return "_data"
}

func (o RepoCrudOptions) UpdateQuery(query *bun.UpdateQuery) *bun.UpdateQuery {
	if o.Columns != nil {
		query = query.Column(o.Columns...)
	}

	return query
}

func (o RepoCrudOptions) BulkUpdateQuery(query *bun.UpdateQuery) *bun.UpdateQuery {
	expr := o.GetExpr()
	query = query.TableExpr(expr)

	if o.Columns != nil {
		for _, col := range o.Columns {
			query = query.Set(fmt.Sprintf("%s = %s.%s", col, expr, col))
		}
	}

	return query
}
