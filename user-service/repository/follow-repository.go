package repository

import (
	"github.com/KristijanPill/Nishtagram/user-service/model"
	"gorm.io/gorm"
)

type FollowRepository struct {
	database *gorm.DB
}

func NewFollowRepository(database *gorm.DB) *FollowRepository {
	return &FollowRepository{database: database}
}

func (repository *FollowRepository) Create(follow *model.Follow) (*model.Follow, error) {
	result := repository.database.Create(follow)

	return follow, result.Error
}

func (repository *FollowRepository) ExistsByUserIDAndFollowerID(userID string, followerID string) bool {
	var follow model.Follow
	return repository.database.Where("user_id = ? AND follower_id = ?", userID, followerID).First(&follow).RowsAffected == 1
}

func (repository *FollowRepository) FindByUserIDAndFollowerID(userID string, followerID string) (*model.Follow, error) {
	var follow model.Follow
	result := repository.database.Preload("Follower").Where("user_id = ? AND follower_id = ?", userID, followerID).First(&follow)

	return &follow, result.Error
}

func (repository *FollowRepository) Update(updatedFollow *model.Follow) (*model.Follow, error) {
	result := repository.database.Model(&model.Follow{}).Where("user_id = ? AND follower_id = ?", updatedFollow.UserID.String(), updatedFollow.FollowerID.String()).Update("close_friend", updatedFollow.CloseFriend).Update("muted", updatedFollow.Muted)

	return updatedFollow, result.Error
}

func (repository *FollowRepository) Delete(follow *model.Follow) error {
	result := repository.database.Delete(follow)

	return result.Error
}

func (repository *FollowRepository) GetFollowerCount(id string) int64 {
	var count int64
	repository.database.Model(&model.Follow{}).Where("user_id = ?", id).Count(&count)

	return count
}

func (repository *FollowRepository) GetFollowingCount(id string) int64 {
	var count int64
	repository.database.Model(&model.Follow{}).Where("follower_id = ?", id).Count(&count)

	return count
}
