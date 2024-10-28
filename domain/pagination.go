package domain

import "strconv"

type Pagination[T any] struct {
	Limit      int    `json:"limit,omitempty" query:"limit"`
	Page       int    `json:"page,omitempty" query:"page"`
	Sort       string `json:"sort,omitempty" query:"sort"`
	TotalRows  int64  `json:"totalRows"`
	TotalPages int    `json:"totalPages"`
	Rows       []T    `json:"rows"`
}

func NewPagination[T any](limit, page, sort string) *Pagination[T] {
	p := &Pagination[T]{}
	p.SetLimit(limit)
	p.SetPage(page)
	p.SetSort(sort)
	return p
}

func (p *Pagination[T]) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination[T]) SetLimit(limit string) {
	if l, err := strconv.Atoi(limit); err == nil {
		p.Limit = l
	} else {
		p.Limit = 10
	}
}

func (p *Pagination[T]) SetPage(page string) {
	if pg, err := strconv.Atoi(page); err == nil {
		p.Page = pg
	} else {
		p.Page = 1
	}
}

func (p *Pagination[T]) SetSort(sort string) {
	if sort != "" {
		p.Sort = sort
	}
}

func (p *Pagination[T]) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination[T]) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination[T]) GetSort() string {
	if p.Sort == "" {
		p.Sort = "Id desc"
	}
	return p.Sort
}

func Map[T, U any](p *Pagination[T], mapper func(T) U) *Pagination[U] {
	newRows := make([]U, len(p.Rows))
	for i, row := range p.Rows {
		newRows[i] = mapper(row)
	}

	return &Pagination[U]{
		Limit:      p.Limit,
		Page:       p.Page,
		Sort:       p.Sort,
		TotalRows:  p.TotalRows,
		TotalPages: p.TotalPages,
		Rows:       newRows,
	}
}
