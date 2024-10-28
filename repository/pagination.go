package repository

import (
	"math"

	"github.com/G-Villarinho/social-network/domain"
	"gorm.io/gorm"
)

func paginate[T any](value *[]T, pagination *domain.Pagination[T], db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64

	db.Model(value).Count(&totalRows)

	pagination.TotalRows = totalRows

	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.GetLimit())))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}
