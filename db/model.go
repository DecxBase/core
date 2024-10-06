package db

import "github.com/DecxBase/core/types"

type CoreModel struct {
}

func (c *CoreModel) KeywordFields() []string {
	return []string{"name"}
}

func (c *CoreModel) UsePagination() bool {
	return true
}

func (c *CoreModel) PaginatePerPage() int32 {
	return 10
}

func (c *CoreModel) ServiceActionToFilters(action types.ServiceActionType) []string {
	return nil
}
