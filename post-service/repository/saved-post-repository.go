package repository

import (
	"github.com/KristijanPill/Nishtagram/post-service/model"
	"gorm.io/gorm"
)

type SavedPostRepository struct {
	database *gorm.DB
}

func NewSavedPostRepository(database *gorm.DB) *SavedPostRepository {
	return &SavedPostRepository{database: database}
}

func (repository *SavedPostRepository) Create(savedPost *model.SavedPost) (*model.SavedPost, error) {
	result := repository.database.Create(savedPost)

	return savedPost, result.Error
}

func (repository *SavedPostRepository) FindAllCollectionNames(userID string) ([]string, error) {
	var collections []string
	result := repository.database.Model(&model.SavedPost{}).Select("distinct(collection_name)").Where("user_id = ?", userID).Find(&collections)

	return collections, result.Error
}

func (repository *SavedPostRepository) FindAllByUserID(userID string) ([]model.SavedPost, error) {
	var saved []model.SavedPost
	result := repository.database.Preload("Post").Where("user_id = ?", userID).Order("created_at desc").Find(&saved)

	return saved, result.Error
}

func (repository *SavedPostRepository) ExistsByPostIDAndUserID(postID string, userID string) (bool, error) {
	var savedPost model.SavedPost
	result := repository.database.Where("post_id = ? AND user_id = ?", postID, userID).Find(&savedPost)

	return result.RowsAffected != 0, result.Error
}

func (repository *SavedPostRepository) Update(savedPost *model.SavedPost) (*model.SavedPost, error) {
	result := repository.database.Save(savedPost)

	return savedPost, result.Error
}

func (repository *SavedPostRepository) Delete(savedPost *model.SavedPost) error {
	result := repository.database.Delete(savedPost)

	return result.Error
}
