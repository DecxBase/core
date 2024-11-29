package repo

import (
	"context"
	"database/sql"

	"github.com/DecxBase/core/types"
)

type RecordsPagingData struct {
	Page    int
	Pages   int
	PerPage int
	Count   int
	Total   int
}

func (r DataRepository[P, M]) Count(ctx context.Context) (int, error) {
	return r.getDB().NewSelect().Model(r.CreateModel()).
		ApplyQueryBuilder(r.ApplyBuilder).
		Count(ctx)
}

func (r DataRepository[P, M]) FindAll(ctx context.Context, data *types.RepoSchemaData, records *[]*P) (*RecordsPagingData, error) {
	total, err := r.KeepFilters().Count(ctx)

	if err != nil {
		return nil, err
	}

	err = r.ApplySelectBuilder(r.ResolveSchemaPaginate(r.getDB().NewSelect().Model(r.CreateModels()), data)).
		ApplyQueryBuilder(r.ApplyBuilder).Scan(ctx, records)

	if err != nil {
		return nil, err
	}

	pages := 0
	page, perPage := r.ExtractPaginationData(data)

	if total > 0 {
		if perPage > 0 {
			pages = total / perPage

			if total > (pages * perPage) {
				pages += 1
			}
		} else {
			page = 1
			pages = 1
			perPage = total
		}
	}

	return &RecordsPagingData{
		Page:    page,
		Pages:   pages,
		PerPage: perPage,
		Count:   len(*records),
		Total:   total,
	}, nil
}

func (r DataRepository[P, M]) Find(ctx context.Context, record *P) error {
	return r.ApplySelectBuilder(r.getDB().NewSelect().Model(r.CreateModel())).
		ApplyQueryBuilder(r.ApplyBuilder).Limit(1).Scan(ctx, record)
}

func (r DataRepository[P, M]) Insert(ctx context.Context, record *P) (sql.Result, error) {
	ref, err := r.Transform(record)

	if err != nil {
		return nil, err
	}

	return r.getDB().NewInsert().Model(ref).Exec(ctx)
}

func (r DataRepository[P, M]) Update(ctx context.Context, record *P, opts RepoCrudOptions) (sql.Result, error) {
	ref, err := r.Transform(record)

	if err != nil {
		return nil, err
	}

	return opts.UpdateQuery(
		r.getDB().NewUpdate().Model(ref).
			OmitZero().
			ApplyQueryBuilder(r.ApplyBuilder),
	).Exec(ctx)
}

func (r DataRepository[P, M]) BulkUpdate(ctx context.Context, records []*P, opts RepoCrudOptions) (sql.Result, error) {
	db := r.getDB()
	values := db.NewValues(&records)

	res, err := opts.BulkUpdateQuery(
		db.NewUpdate().With(opts.GetExpr(), values).Model(r.CreateModel()).OmitZero(),
	).ApplyQueryBuilder(r.ApplyBuilder).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r DataRepository[P, M]) Delete(ctx context.Context, records ...*P) (sql.Result, error) {
	query := r.getDB().NewDelete()

	if len(records) > 0 {
		rlist, err := r.TransformMulti(records)

		if err != nil {
			return nil, err
		}

		query = query.Model(&rlist).WherePK()
	} else {
		query = query.Model(r.CreateModel())
	}

	res, err := query.ApplyQueryBuilder(r.ApplyBuilder).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return res, nil
}
