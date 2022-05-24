package repository

import (
	"github.com/KristijanPill/Nishtagram/user-service/model"
	"gorm.io/gorm"
)

type FollowRequestRepository struct {
	database *gorm.DB
}

func NewFollowRequestRepository(database *gorm.DB) *FollowRequestRepository {
	return &FollowRequestRepository{database: database}
}

func (repository *FollowRequestRepository) Create(request *model.FollowRequest) (*model.FollowRequest, error) {
	result := repository.database.Create(request)

	return request, result.Error
}

func (repository *FollowRequestRepository) FindByUserIDAndFollowedID(userID string, followedID string) (*model.FollowRequest, error) {
	var request model.FollowRequest
	result := repository.database.Where("user_id = ? AND followed_id = ?", userID, followedID).First(&request)

	return &request, result.Error
}

func (repository *FollowRequestRepository) ExistsByUserIDAndFollowedID(userID string, followedID string) bool {
	var request model.FollowRequest
	return repository.database.Where("user_id = ? AND followed_id = ?", userID, followedID).First(&request).RowsAffected == 1
}

func (repository *FollowRequestRepository) Delete(request *model.FollowRequest) error {
	result := repository.database.Delete(request)

	return result.Error
}
