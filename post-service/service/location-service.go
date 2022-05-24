package service

import (
	"github.com/KristijanPill/Nishtagram/post-service/model"
	"github.com/KristijanPill/Nishtagram/post-service/repository"
)

type LocationService struct {
	repository *repository.LocationRepository
}

func NewLocationService(repository *repository.LocationRepository) *LocationService {
	return &LocationService{repository: repository}
}

func (service *LocationService) GetByQuery(query string) ([]model.Location, error) {
	return service.repository.GetByQuery(query)
}
