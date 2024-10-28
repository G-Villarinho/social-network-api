package repository

import (
	"math"

	"github.com/G-Villarinho/social-network/domain"
	"gorm.io/gorm"
)

func paginate[T any](pagination *domain.Pagination[T], db *gorm.DB) (*domain.Pagination[T], error) {
	var totalRows int64

	if err := db.Model(&pagination.Rows).Count(&totalRows).Error; err != nil {
		return nil, err
	}
	pagination.TotalRows = totalRows

	pagination.TotalPages = int(math.Ceil(float64(totalRows) / float64(pagination.GetLimit())))

	if err := db.Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit()).
		Order(pagination.GetSort()).
		Find(&pagination.Rows).Error; err != nil {
		return nil, err
	}

	return pagination, nil
}
