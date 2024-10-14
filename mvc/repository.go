package mvc

import (
	"wgenerator/definitions"

	"gorm.io/gorm"
)

type QueryCondition struct {
	Query string
	Args  []interface{}
}

type BaseRepository[T any] struct {
	DBF       *gorm.DB
	TableName string
}

func (br BaseRepository[T]) ReadlyTableQuery() *gorm.DB {
	return br.DBF.Table(br.TableName).Debug()
}

func (br BaseRepository[T]) GetById(id int) (T, error) {
	return br.GetOneWithCondition(QueryCondition{"id = ?", []interface{}{id}})
}

func (br BaseRepository[T]) GetOneWithCondition(qc QueryCondition) (T, error) {
	var model T

	if error := br.ReadlyTableQuery().Where(qc.Query, qc.Args...).Limit(1).First(&model).Error; error != nil {
		return model, error
	}
	return model, nil
}

func (br BaseRepository[T]) Count(qc QueryCondition) (int64, error) {
	var count int64

	if error := br.ReadlyTableQuery().Where(qc.Query, qc.Args...).Count(&count).Error; error != nil {
		return 0, error
	}
	return count, nil
}

func (br BaseRepository[T]) IsExists(qc QueryCondition) (bool, error) {
	count, error := br.Count(qc)
	return count > 0, error
}

func (br BaseRepository[T]) ListByIds(ids []int) ([]T, error) {
	return br.ListWithCondition(QueryCondition{"id IN ?", []interface{}{ids}})
}

func (br BaseRepository[T]) ListWithCondition(qc QueryCondition) ([]T, error) {
	var models []T

	if error := br.ReadlyTableQuery().Where(qc.Query, qc.Args...).Find(&models).Error; error != nil {
		return models, error
	}
	return models, nil
}

func (br BaseRepository[T]) ListWithConditionByPagination(qc QueryCondition, page, pageSize int) (definitions.Pagination[T], error) {
	var pagination definitions.Pagination[T]

	if error := br.ReadlyTableQuery().Where(qc.Query, qc.Args...).Count(&pagination.Total).Offset((page - 1) * pageSize).Limit(pageSize).Find(&pagination.Datas).Error; error != nil {
		return pagination, error
	}
	return pagination, nil
}

func (br BaseRepository[T]) RemoveByIds(ids []int) (int64, error) {
	return br.RemoveWithCondition(QueryCondition{"id IN ?", []interface{}{ids}})
}

func (br BaseRepository[T]) RemoveWithCondition(qc QueryCondition) (int64, error) {
	var models []T

	result := br.ReadlyTableQuery().Where(qc.Query, qc.Args...).Delete(&models)

	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (br BaseRepository[T]) SaveBatch(models []T, batchSize int) (int64, error) {
	result := br.ReadlyTableQuery().CreateInBatches(models, batchSize)

	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (br BaseRepository[T]) Save(model T) (int64, error) {
	result := br.ReadlyTableQuery().Create(&model)

	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (br BaseRepository[T]) SaveOrUpdate(model T) (int64, error) {
	result := br.ReadlyTableQuery().Save(&model)

	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (br BaseRepository[T]) UpdateColumns(qc QueryCondition, model T) (int64, error) {
	var tm T

	result := br.ReadlyTableQuery().Model(&tm).Where(qc.Query, qc.Args...).Updates(model)

	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

type IRepository[T any] interface {
	Count(qc QueryCondition) (int64, error)
	IsExists(qc QueryCondition) (bool, error)
	GetById(id int) (T, error)
	GetOneWithCondition(qc QueryCondition) (T, error)
	ListWithCondition(qc QueryCondition) ([]T, error)
	ListWithConditionByPagination(qc QueryCondition, page, pageSize int) (definitions.Pagination[T], error)
	RemoveByIds(ids []int) (int64, error)
	RemoveWithCondition(qc QueryCondition) (int64, error)
	SaveBatch(models []T, batchSize int) (int64, error)
	Save(model T) (int64, error)
	SaveOrUpdate(model T) (int64, error)
	UpdateColumns(qc QueryCondition, model T) (int64, error)
}
