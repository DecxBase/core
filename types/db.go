package types

type DataFilterCondition int

const (
	DataFilterWhere     DataFilterCondition = iota + 1 // EnumIndex = 1
	DataFilterWherePK                                  // EnumIndex = 2
	DataFilterWhereLike                                // EnumIndex = 3
	DataFilterFunc                                     // EnumIndex = 4
)

func (r DataFilterCondition) String() string {
	return []string{"where", "pk", "like", "func"}[r-1]
}

type ServiceActionType int

const (
	ServiceActionFindAll ServiceActionType = iota + 1 // EnumIndex = 1
	ServiceActionFind                                 // EnumIndex = 2
	ServiceActionDelete                               // EnumIndex = 3
)

func (r ServiceActionType) String() string {
	return []string{"find_all", "find", "delete"}[r-1]
}

type ModelForService interface {
	ModelForService() string
	ServiceActionToFilters(ServiceActionType) []string
	KeywordFields() []string
	UsePagination() bool
	PaginatePerPage() int32
}

type RepoSchemaData struct {
	Filter   map[string]any
	Paginate map[string]any
}
