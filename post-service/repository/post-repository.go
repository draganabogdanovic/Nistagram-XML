package repository

import (
	"strings"

	"github.com/KristijanPill/Nishtagram/post-service/model"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type PostRepository struct {
	database *gorm.DB
}

func NewPostRepository(database *gorm.DB) *PostRepository {
	return &PostRepository{database: database}
}

func (repository *PostRepository) Create(post *model.Post) (*model.Post, error) {
	result := repository.database.Create(post)

	return post, result.Error
}

func (repository *PostRepository) FindById(id string) (*model.Post, error) {
	var post model.Post
	result := repository.database.Preload("Location").First(&post, "id = ?", id)

	return &post, result.Error
}

func (repository *PostRepository) FindByUserID(userID string) ([]model.Post, error) {
	var posts []model.Post
	result := repository.database.Preload("Location").Where("user_id = ? ", userID).Order("created_at desc").Find(&posts)

	return posts, result.Error
}

func (repository *PostRepository) FindLikeLocation(query string) ([]model.Post, error) {
	var posts []model.Post
	result := repository.database.Joins("Location").Where("country ILIKE ?", "%"+query+"%").Or("city ILIKE ?", "%"+query+"%").Find(&posts)

	return posts, result.Error
}

func (repository *PostRepository) FindLikeTags(query string) ([]model.Post, error) {
	var posts []model.Post
	tags := strings.Split(query, " ")
	result := repository.database.Where("tags && ?", pq.Array(tags)).Find(&posts)

	return posts, result.Error
}

func (repository *PostRepository) CreateReport(report *model.Report) (*model.Report, error) {
	result := repository.database.Create(report)

	return report, result.Error
}
