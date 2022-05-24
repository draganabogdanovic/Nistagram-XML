package repository

import (
	"github.com/KristijanPill/Nishtagram/post-service/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReviewRepository struct {
	database *gorm.DB
}

func NewReviewRepository(database *gorm.DB) *ReviewRepository {
	return &ReviewRepository{database: database}
}

func (repository *ReviewRepository) Create(review *model.Review) (*model.Review, error) {
	result := repository.database.Create(review)

	return review, result.Error
}

func (repository *ReviewRepository) ExistsByPostIDAndUserID(postID string, userID string) (bool, error) {
	var review model.Review
	result := repository.database.Where("post_id = ? AND user_id = ?", postID, userID).Find(&review)

	return result.RowsAffected != 0, result.Error
}

func (repository *ReviewRepository) Update(review *model.Review) (*model.Review, error) {
	result := repository.database.Save(review)

	return review, result.Error
}

func (repository *ReviewRepository) Delete(postID string, userID string) error {
	result := repository.database.Where("post_id = ? AND user_id = ?", postID, userID).Delete(&model.Review{})

	return result.Error
}

func (repository *ReviewRepository) FindCountByPostIDAndStatus(postID string, status model.ReviewStatus) int64 {
	var reviews []model.Review
	result := repository.database.Where("post_id = ? AND status = ?", postID, status).Find(&reviews)

	return result.RowsAffected
}

func (repository *ReviewRepository) FindStatusByPostIDAndUserID(postID string, userID string) int {
	var review model.Review
	result := repository.database.Where("post_id = ? AND user_id = ?", postID, userID).Find(&review)

	if result.RowsAffected == 0 {
		return 2
	} else {
		return int(review.Status)
	}
}

func (repository *ReviewRepository) GetReviewsByUserIDAndStatus(userID uuid.UUID, status int) ([]model.Review, error) {
	var reviews []model.Review
	result := repository.database.Preload("Post").Where("user_id = ? AND status = ?", userID, status).Find(&reviews)
	return reviews, result.Error
}
