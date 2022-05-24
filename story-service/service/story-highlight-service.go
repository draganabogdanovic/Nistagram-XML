package service

import (
	"errors"

	"github.com/KristijanPill/Nishtagram/story-service/model"
	"github.com/KristijanPill/Nishtagram/story-service/payload"
	"github.com/KristijanPill/Nishtagram/story-service/repository"
	"github.com/google/uuid"
)

type StoryHighlightService struct {
	highlightRepository *repository.StoryHighlightRepository
	storyRepository     *repository.StoryRepository
}

func NewStoryHighlightService(highlightRepository *repository.StoryHighlightRepository, storyRepository *repository.StoryRepository) *StoryHighlightService {
	return &StoryHighlightService{
		highlightRepository: highlightRepository,
		storyRepository:     storyRepository,
	}
}

func (service *StoryHighlightService) HighlightStory(dto *payload.StoryHighlightCreate) (*model.StoryHighlight, error) {
	story, err := service.storyRepository.FindByID(dto.StoryID.String())

	if err != nil {
		return nil, err
	}

	if story.UserID != dto.UserID {
		return nil, errors.New("unauthorized")
	}

	exists, err := service.highlightRepository.ExistsByPostIDAndUserID(dto.StoryID.String(), dto.UserID.String())

	if err != nil {
		return nil, err
	}

	highlight := &model.StoryHighlight{
		UserID:        dto.UserID,
		StoryID:       dto.StoryID,
		HighlightName: dto.HighlightName,
	}

	if exists {
		return service.highlightRepository.Update(highlight)
	} else {
		return service.highlightRepository.Create(highlight)
	}
}

func (service *StoryHighlightService) GetAllHighlightNames(loggedInUserID uuid.UUID) ([]string, error) {
	return service.highlightRepository.FindAllHighlightNames(loggedInUserID.String())
}

func (service *StoryHighlightService) GetAllByLoggedInUser(loggedInUserID uuid.UUID) ([]payload.StoryHighlightView, error) {
	highlights, err := service.highlightRepository.FindAllByUserID(loggedInUserID.String())

	if err != nil {
		return nil, err
	}

	var highlightsView []payload.StoryHighlightView

	for _, highlight := range highlights {
		storyView := &payload.StoryView{
			ID:               highlight.Story.ID,
			CreatedAt:        highlight.Story.CreatedAt,
			Content:          highlight.Story.Content,
			CloseFriendsOnly: highlight.Story.CloseFriendsOnly,
		}
		highlightView := &payload.StoryHighlightView{
			StoryView:     *storyView,
			HighlightName: highlight.HighlightName,
		}
		highlightsView = append(highlightsView, *highlightView)
	}

	return highlightsView, nil
}
