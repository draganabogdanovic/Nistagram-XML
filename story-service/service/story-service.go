package service

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/KristijanPill/Nishtagram/story-service/helpers"
	"github.com/KristijanPill/Nishtagram/story-service/model"
	"github.com/KristijanPill/Nishtagram/story-service/payload"
	"github.com/KristijanPill/Nishtagram/story-service/repository"
	"github.com/google/uuid"
)

type StoryService struct {
	repository *repository.StoryRepository
}

func NewStoryService(repository *repository.StoryRepository) *StoryService {
	return &StoryService{repository: repository}
}

func (service *StoryService) Create(story *model.Story) (*model.Story, error) {
	return service.repository.Create(story)
}

func (service *StoryService) Delete(storyID uuid.UUID, userID uuid.UUID) error {
	story, err := service.repository.FindByID(storyID.String())

	if err != nil {
		return err
	}

	if story.UserID == userID {
		return service.repository.Delete(story)
	} else {
		return errors.New("unauthorized")
	}
}

func (service *StoryService) FindByUser(userID uuid.UUID, loggedInUserID uuid.UUID, token string) ([]payload.StoryView, error) {
	requestURL := fmt.Sprintf("http://%s:%s/follow/outgoing/"+userID.String(), os.Getenv("USER_SERVICE_DOMAIN"), os.Getenv("USER_SERVICE_PORT"))
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, requestURL, nil)

	if token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	}

	response, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	var followStatus = &payload.FollowStatus{}
	helpers.FromJSON(&followStatus, response.Body)

	var stories = []model.Story{}
	if followStatus.CloseFriend {
		stories, err = service.repository.FindByUserIDCloseFriends(userID.String())
	} else {
		stories, err = service.repository.FindByUserIDNotCloseFriends(userID.String())
	}

	if err != nil {
		return nil, err
	}

	var storiesView []payload.StoryView

	for _, story := range stories {
		storyView := &payload.StoryView{
			ID:               story.ID,
			CreatedAt:        story.CreatedAt,
			Content:          story.Content,
			CloseFriendsOnly: story.CloseFriendsOnly,
		}
		storiesView = append(storiesView, *storyView)
	}

	return storiesView, nil
}

func (service *StoryService) FindByLoggedInUser(userID uuid.UUID) ([]payload.StoryView, error) {

	stories, err := service.repository.FindByUserID(userID.String())

	if err != nil {
		return nil, err
	}

	var storiesView []payload.StoryView

	for _, story := range stories {
		storyView := &payload.StoryView{
			ID:               story.ID,
			CreatedAt:        story.CreatedAt,
			Content:          story.Content,
			CloseFriendsOnly: story.CloseFriendsOnly,
		}
		storiesView = append(storiesView, *storyView)
	}

	return storiesView, nil
}

func (service *StoryService) FindAllByLoggedInUser(userID uuid.UUID) ([]payload.StoryView, error) {
	stories, err := service.repository.FindAllByUserID(userID.String())

	if err != nil {
		return nil, err
	}

	var storiesView []payload.StoryView

	for _, story := range stories {
		storyView := &payload.StoryView{
			ID:               story.ID,
			CreatedAt:        story.CreatedAt,
			Content:          story.Content,
			CloseFriendsOnly: story.CloseFriendsOnly,
		}
		storiesView = append(storiesView, *storyView)
	}

	return storiesView, nil
}

func (service *StoryService) CreateReport(dto *payload.ReportCreate) (*model.Report, error) {

	report := &model.Report{
		StoryID:     dto.StoryID,
		UserID:      dto.UserID,
		Description: dto.Description,
	}

	return service.repository.CreateReport(report)
}
