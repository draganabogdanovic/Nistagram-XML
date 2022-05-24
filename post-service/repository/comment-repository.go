package repository

import (
	"github.com/KristijanPill/Nishtagram/post-service/model"
	"gorm.io/gorm"
)

type CommentRepository struct {
	database *gorm.DB
}

func NewCommentRepository(database *gorm.DB) *CommentRepository {
	return &CommentRepository{database: database}
}

func (repository *CommentRepository) Create(comment *model.Comment) (*model.Comment, error) {
	result := repository.database.Create(comment)

	return comment, result.Error
}

func (repository *CommentRepository) FindById(id string) (*model.Comment, error) {
	var comment model.Comment
	result := repository.database.First(&comment, "id = ?", id)

	return &comment, result.Error
}

func (repository *CommentRepository) FindAllByPostID(commentID string) ([]model.Comment, error) {
	var comments []model.Comment
	result := repository.database.Where("post_id = ?", commentID).Find(&comments)

	return comments, result.Error
}

func (repository *CommentRepository) Delete(postID string) error {
	result := repository.database.Delete(&model.Comment{}, postID)

	return result.Error
}

func (repository *CommentRepository) FindCountByPostID(postID string) int64 {
	var comments []model.Comment
	result := repository.database.Where("post_id = ?", postID).Find(&comments)

	return result.RowsAffected
}
