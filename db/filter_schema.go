package db

import (
	"fmt"

	"github.com/uptrace/bun"
)

type FilterSchemas = map[string]FilterSchema

type FilterSchema struct {
	Translate string
	Condition string
	Multiples []string
	Builder   func(bun.QueryBuilder, any) bun.QueryBuilder
}

func (s FilterSchema) Apply(q bun.QueryBuilder, key string, value any) bun.QueryBuilder {
	if s.Builder != nil {
		return s.Builder(q, value)
	}

	if s.Multiples != nil {

	}
	if len(s.Translate) > 0 {
		key = s.Translate
	}

	return ApplyCondition(q, s.Condition, key, value)
}

func ApplyCondition(q bun.QueryBuilder, condition string, key string, value any) bun.QueryBuilder {
	switch condition {
	case "like":
		q = ApplyLike(q, key, value)
	default:
		q = q.Where(fmt.Sprintf("%s = ?", key), value)
	}

	return q
}

func ApplyEqual(q bun.QueryBuilder, key string, value any, useOr ...bool) bun.QueryBuilder {
	useAnd := true
	if len(useOr) > 0 {
		useAnd = !useOr[0]
	}

	if useAnd
	return q.Where(fmt.Sprintf("%s = ?", key), value)
}

func ApplyLike(q bun.QueryBuilder, key string, value any) bun.QueryBuilder {
	return q.Where(fmt.Sprintf("%s LIKE ?", key), "%"+fmt.Sprintf("%v", value)+"%")
}
