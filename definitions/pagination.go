package definitions

type Pagination[T any] struct {
	Total int64 `json:"total"`
	Datas []T   `json:"datas"`
}

func NewPagination[T any](total int64, datas []T) Pagination[T] {
	return Pagination[T]{
		Total: total, Datas: datas,
	}
}
