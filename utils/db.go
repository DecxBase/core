package utils

import (
	"fmt"

	"github.com/DecxBase/core/types"
	"github.com/uptrace/bun"
)

func ApplyCondition(
	q bun.QueryBuilder, condition types.DataFilterCondition,
	useOR bool, key string, values ...any,
) bun.QueryBuilder {
	switch condition {
	case types.DataFilterWherePK:
		q = q.WherePK()
	case types.DataFilterWhereLike:
		q = ApplyQueryLike(q, useOR, key, values[0])
	default:
		q = ApplyQueryEqual(q, useOR, key, values...)
	}

	return q
}

func ApplyQueryWhere(q bun.QueryBuilder, useOR bool, key string, values ...any) bun.QueryBuilder {
	if useOR {
		return q.WhereOr(key, values...)
	}

	return q.Where(key, values...)
}

func ApplyQueryEqual(q bun.QueryBuilder, useOR bool, key string, values ...any) bun.QueryBuilder {
	return ApplyQueryWhere(q, useOR, fmt.Sprintf("%s = ?", key), values...)
}

func ApplyQueryLike(q bun.QueryBuilder, useOR bool, key string, value any) bun.QueryBuilder {
	return ApplyQueryWhere(q, useOR, fmt.Sprintf("%s LIKE ?", key), "%"+fmt.Sprintf("%v", value)+"%")
}

func ParseSchemaData(action types.ServiceActionType, source any) *types.RepoSchemaData {
	resolved, err := StructToMap(source)
	if err != nil {
		return nil
	}

	var filters map[string]any
	fData := resolved["filter"]
	if fData != nil {
		filters = fData.(map[string]any)
	}

	var paginate map[string]any
	if action == types.ServiceActionFindAll {
		pData := resolved["paginate"]
		if pData != nil {
			paginate = pData.(map[string]any)
		}
	}

	if filters == nil && paginate == nil && len(resolved) > 0 {
		filters = resolved
	}

	return &types.RepoSchemaData{
		Filter:   filters,
		Paginate: paginate,
	}
}
