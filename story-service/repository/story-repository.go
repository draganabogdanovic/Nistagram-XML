package repository

import (
	"time"

	"github.com/KristijanPill/Nishtagram/story-service/model"
	"gorm.io/gorm"
)

type StoryRepository struct {
	database *gorm.DB
}

func NewStoryRepository(database *gorm.DB) *StoryRepository {
	return &StoryRepository{database: database}
}

func (repository *StoryRepository) Create(story *model.Story) (*model.Story, error) {
	result := repository.database.Create(story)

	return story, result.Error
}

func (repository *StoryRepository) FindByID(storyID string) (*model.Story, error) {
	var story model.Story
	result := repository.database.First(&story, "id = ?", storyID)

	return &story, result.Error
}

func (repository *StoryRepository) FindByUserIDNotCloseFriends(userID string) ([]model.Story, error) {
	var stories []model.Story
	result := repository.database.Where("user_id = ? AND close_friends_only = ? AND created_at >= ?", userID, false, time.Now().Add(-24*time.Hour)).Order("created_at desc").Find(&stories)

	return stories, result.Error
}

func (repository *StoryRepository) FindByUserIDCloseFriends(userID string) ([]model.Story, error) {
	var stories []model.Story
	result := repository.database.Where("user_id = ? AND created_at >= ?", userID, time.Now().Add(-24*time.Hour)).Order("created_at desc").Find(&stories)

	return stories, result.Error
}

func (repository *StoryRepository) FindByUserID(userID string) ([]model.Story, error) {
	var stories []model.Story
	result := repository.database.Where("user_id = ? AND created_at >= ?", userID, time.Now().Add(-24*time.Hour)).Order("created_at desc").Find(&stories)

	return stories, result.Error
}

func (repository *StoryRepository) FindAllByUserID(userID string) ([]model.Story, error) {
	var stories []model.Story
	result := repository.database.Where("user_id = ?", userID).Order("created_at desc").Find(&stories)

	return stories, result.Error
}

func (repository *StoryRepository) Delete(story *model.Story) error {
	result := repository.database.Delete(story)

	return result.Error
}

func (repository *StoryRepository) CreateReport(report *model.Report) (*model.Report, error) {
	result := repository.database.Create(report)

	return report, result.Error
}
