package mvc

import "github.com/cabbagen/wgenerator/definitions"

type BaseService[T any] struct {
	Repositoriy IRepository[T]
}

func (tbs BaseService[T]) Count(qc QueryCondition) (int64, error) {
	return tbs.Repositoriy.Count(qc)
}

func (tbs BaseService[T]) IsExists(qc QueryCondition) (bool, error) {
	return tbs.Repositoriy.IsExists(qc)
}

func (tbs BaseService[T]) GetById(id int) (T, error) {
	return tbs.Repositoriy.GetById(id)
}

func (tbs BaseService[T]) GetOneWithCondition(qc QueryCondition) (T, error) {
	return tbs.Repositoriy.GetOneWithCondition(qc)
}

func (tbs BaseService[T]) ListWithCondition(qc QueryCondition) ([]T, error) {
	return tbs.Repositoriy.ListWithCondition(qc)
}

func (tbs BaseService[T]) ListWithConditionByPagination(qc QueryCondition, page, pageSize int) (definitions.Pagination[T], error) {
	return tbs.Repositoriy.ListWithConditionByPagination(qc, page, pageSize)
}

func (tbs BaseService[T]) RemoveByIds(ids []int) (int64, error) {
	return tbs.Repositoriy.RemoveByIds(ids)
}

func (tbs BaseService[T]) RemoveWithCondition(qc QueryCondition) (int64, error) {
	return tbs.Repositoriy.RemoveWithCondition(qc)
}

func (tbs BaseService[T]) SaveBatch(models []T, batchSize int) (int64, error) {
	return tbs.Repositoriy.SaveBatch(models, batchSize)
}

func (tbs BaseService[T]) Save(model T) (int64, error) {
	return tbs.Repositoriy.Save(model)
}

func (tbs BaseService[T]) SaveOrUpdate(model T) (int64, error) {
	return tbs.Repositoriy.SaveOrUpdate(model)
}

func (tbs BaseService[T]) UpdateColumns(qc QueryCondition, model T) (int64, error) {
	return tbs.Repositoriy.UpdateColumns(qc, model)
}
