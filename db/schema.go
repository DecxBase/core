package db

func (r *DataRepository[P, M]) RegisterSchema(schema FilterSchemas) {
	r.schema = schema
}
