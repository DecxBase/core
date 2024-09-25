package db

import (
	"encoding/json"
	"reflect"

	"github.com/uptrace/bun"
)

func (r *DataRepository[P, M]) getDB() *bun.DB {
	if r.tempDB != nil {
		db := r.tempDB
		r.tempDB = nil

		return db
	}

	return r.db
}

func (r *DataRepository[P, M]) Dispose() {
	if r.tempDB != nil {
		r.tempDB.Close()
		r.tempDB = nil
	}

	r.db.Close()
}

func (r *DataRepository[P, M]) WithDB(db *bun.DB) *DataRepository[P, M] {
	r.tempDB = db
	return r
}

func (r *DataRepository[P, M]) WithFilters(filters ...*QueryFilters) *DataRepository[P, M] {
	r.activeFilters = filters
	return r
}

func (r *DataRepository[P, M]) AddFilter(filter *QueryFilters) *DataRepository[P, M] {
	r.activeFilters = append(r.activeFilters, filter)
	return r
}

func (r *DataRepository[P, M]) ApplyPaginate(query *bun.SelectQuery) *bun.SelectQuery {
	if len(r.activeFilters) > 0 {
		for _, filter := range r.activeFilters {
			query = filter.ApplyPaginate(query)
		}
	}

	return query
}

func (r *DataRepository[P, M]) ApplyBuilder(query bun.QueryBuilder) bun.QueryBuilder {
	if len(r.activeFilters) > 0 {
		for _, filter := range r.activeFilters {
			query = filter.ApplyBuilder(query)
		}
	}

	r.activeFilters = make([]*QueryFilters, 0)

	return query
}

func (r DataRepository[P, M]) CreateProto() *P {
	var zero *P
	tp := reflect.TypeOf(zero).Elem()
	ref := reflect.New(tp).Interface()

	return ref.(*P)
}

func (r DataRepository[P, M]) CreateProtos() *[]*P {
	var zero [0]*P
	tp := reflect.TypeOf(zero).Elem()
	ref := reflect.New(reflect.SliceOf(tp)).Interface()

	return ref.(*[]*P)
}

func (r DataRepository[P, M]) CreateModel() *M {
	var zero *M
	tp := reflect.TypeOf(zero).Elem()
	ref := reflect.New(tp).Interface()

	return ref.(*M)
}

func (r DataRepository[P, M]) CreateModels() *[]*M {
	var zero [0]*M
	tp := reflect.TypeOf(zero).Elem()
	ref := reflect.New(reflect.SliceOf(tp)).Interface()

	return ref.(*[]*M)
}

func (r DataRepository[P, M]) Transform(record *P) (*M, error) {
	ref := r.CreateModel()

	enc, err := json.Marshal(record)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(enc, ref)
	if err != nil {
		return nil, err
	}

	return ref, nil
}

func (r DataRepository[P, M]) TransformMulti(records []*P) ([]*M, error) {
	var err error
	var record *M
	list := make([]*M, 0)

	for _, item := range records {
		record, err = r.Transform(item)

		if err != nil {
			break
		} else {
			list = append(list, record)
		}
	}

	return list, err
}

func (r DataRepository[P, M]) RTransform(record *M) (*P, error) {
	ref := r.CreateProto()

	enc, err := json.Marshal(record)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(enc, ref)
	if err != nil {
		return nil, err
	}

	return ref, nil
}

func (r DataRepository[P, M]) RTransformMulti(records []*M) ([]*P, error) {
	var err error
	var record *P
	list := make([]*P, 0)

	for _, item := range records {
		record, err = r.RTransform(item)

		if err != nil {
			break
		} else {
			list = append(list, record)
		}
	}

	return list, err
}
