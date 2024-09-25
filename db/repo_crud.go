package db

import (
	"context"
	"database/sql"
)

func (r DataRepository[P, M]) Find(ctx context.Context, record *P) error {
	return r.getDB().NewSelect().Model(r.CreateProto()).
		ApplyQueryBuilder(r.ApplyBuilder).Limit(1).
		Scan(ctx, record)
}

func (r DataRepository[P, M]) FindAll(ctx context.Context, records *[]*P, recordsRef *[]*M) error {
	return r.ApplyPaginate(r.getDB().NewSelect().Model(r.CreateProtos())).
		ApplyQueryBuilder(r.ApplyBuilder).
		Scan(ctx, records)
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
