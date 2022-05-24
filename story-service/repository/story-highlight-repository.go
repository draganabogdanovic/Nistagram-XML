package repository

import (
	"github.com/KristijanPill/Nishtagram/story-service/model"
	"gorm.io/gorm"
)

type StoryHighlightRepository struct {
	database *gorm.DB
}

func NewStoryHighlightRepository(database *gorm.DB) *StoryHighlightRepository {
	return &StoryHighlightRepository{database: database}
}

func (repository *StoryHighlightRepository) FindAllHighlightNames(userID string) ([]string, error) {
	var highlightNames []string
	result := repository.database.Model(&model.StoryHighlight{}).Select("distinct(highlight_name)").Where("user_id = ?", userID).Find(&highlightNames)

	return highlightNames, result.Error
}

func (repository *StoryHighlightRepository) FindAllByUserID(userID string) ([]model.StoryHighlight, error) {
	var highlights []model.StoryHighlight
	result := repository.database.Preload("Story").Where("user_id = ?", userID).Order("created_at desc").Find(&highlights)

	return highlights, result.Error
}

func (repository *StoryHighlightRepository) Create(highlight *model.StoryHighlight) (*model.StoryHighlight, error) {
	result := repository.database.Create(highlight)

	return highlight, result.Error
}

func (repository *StoryHighlightRepository) ExistsByPostIDAndUserID(storyID string, userID string) (bool, error) {
	var highlight model.StoryHighlight
	result := repository.database.Where("story_id = ? AND user_id = ?", storyID, userID).Find(&highlight)

	return result.RowsAffected != 0, result.Error
}

func (repository *StoryHighlightRepository) Update(highlight *model.StoryHighlight) (*model.StoryHighlight, error) {
	result := repository.database.Save(highlight)

	return highlight, result.Error
}

func (repository *StoryHighlightRepository) Delete(highlight *model.StoryHighlight) error {
	result := repository.database.Delete(highlight)

	return result.Error
}
