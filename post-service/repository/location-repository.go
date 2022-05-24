package repository

import (
	"github.com/KristijanPill/Nishtagram/post-service/model"
	"gorm.io/gorm"
)

type LocationRepository struct {
	database *gorm.DB
}

func NewLocationRepository(database *gorm.DB) *LocationRepository {
	return &LocationRepository{database: database}
}

func (repository *LocationRepository) GetByQuery(query string) ([]model.Location, error) {
	var locations []model.Location
	result := repository.database.Limit(5).Where("country ILIKE ?", "%"+query+"%").Or("city ILIKE ?", "%"+query+"%").Find(&locations)

	return locations, result.Error
}
