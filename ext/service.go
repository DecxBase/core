package ext

import (
	"context"

	"github.com/DecxBase/core/exception"
	"github.com/DecxBase/core/repo"
	"github.com/DecxBase/core/types"
	"github.com/uptrace/bun"
)

type Service[P any, M types.ModelForService] struct {
	Repo *repo.DataRepository[P, M]
}

func NewService[P any, M types.ModelForService](
	ins *bun.DB,
	configure ...func(*repo.DataRepository[P, M]) *repo.DataRepository[P, M],
) *Service[P, M] {
	repo := repo.NewRepository[P, M](ins)

	if len(configure) > 0 {
		repo = configure[0](repo)
	}

	return &Service[P, M]{
		Repo: repo,
	}
}

func (s *Service[P, M]) FindRecords(ctx context.Context, req any) ([]*P, *repo.RecordsPagingData, error) {
	var records []*P
	data := s.Repo.ResolveSchemaData(types.ServiceActionFindAll, req)
	filter := s.Repo.ResolveSchemaFilter(types.ServiceActionFindAll, data)
	paging, err := s.Repo.WithFilters(filter).FindAll(ctx, data, &records)

	if err != nil {
		return nil, nil, exception.Raise(err).WithCode("find_err")
	}

	return records, paging, nil
}

func (s *Service[P, M]) FindRecord(ctx context.Context, req any) (*P, error) {
	record := new(P)
	data := s.Repo.ResolveSchemaData(types.ServiceActionFind, req)

	err := s.Repo.WithFilters(
		s.Repo.ResolveSchemaFilter(types.ServiceActionFind, data),
	).Find(ctx, record)

	if err != nil {
		return nil, exception.Raise(err).WithCode("find_err").WithMessage("failed to find record")
	}

	return record, nil
}

func (s *Service[P, M]) SaveRecord(ctx context.Context, create bool, record *P) error {
	var err error

	if create {
		_, err = s.Repo.Insert(ctx, record)
	} else {
		filters := repo.NewDataFilters().WherePK()
		_, err = s.Repo.WithFilters(filters).Update(ctx, record, repo.RepoCrudOptions{})
	}

	if err != nil {
		return exception.Raise(err).WithCode("save_err")
	}

	return err
}

func (s *Service[P, M]) DeleteRecord(ctx context.Context, req any) error {
	data := s.Repo.ResolveSchemaData(types.ServiceActionDelete, req)

	_, err := s.Repo.WithFilters(
		s.Repo.ResolveSchemaFilter(types.ServiceActionDelete, data),
	).Delete(ctx)

	if err != nil {
		return exception.Raise(err).WithCode("delete_err")
	}

	return err
}
