package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/KristijanPill/Nishtagram/post-service/helpers"
	"github.com/KristijanPill/Nishtagram/post-service/model"
	"github.com/KristijanPill/Nishtagram/post-service/payload"
	"github.com/KristijanPill/Nishtagram/post-service/repository"
	"github.com/google/uuid"
)

type SavedPostService struct {
	savedPostRepository *repository.SavedPostRepository
	reviewRepository    *repository.ReviewRepository
	commentRepository   *repository.CommentRepository
}

func NewSavedPostService(savedPostRepository *repository.SavedPostRepository, reviewRepository *repository.ReviewRepository, commentRepository *repository.CommentRepository) *SavedPostService {
	return &SavedPostService{
		savedPostRepository: savedPostRepository,
		reviewRepository:    reviewRepository,
		commentRepository:   commentRepository,
	}
}

func (service *SavedPostService) SavePost(dto *payload.SavedPostCreate) (*model.SavedPost, error) {
	exists, err := service.savedPostRepository.ExistsByPostIDAndUserID(dto.PostID.String(), dto.UserID.String())

	if err != nil {
		return nil, err
	}

	savedPost := &model.SavedPost{
		UserID:         dto.UserID,
		PostID:         dto.PostID,
		CollectionName: dto.CollectionName,
	}

	if exists {
		return service.savedPostRepository.Update(savedPost)
	} else {
		return service.savedPostRepository.Create(savedPost)
	}
}

func (service *SavedPostService) RemoveSavedPost(savedPost *model.SavedPost) error {
	return service.savedPostRepository.Delete(savedPost)
}

func (service *SavedPostService) GetAllCollectionNames(loggedInUserID uuid.UUID) ([]string, error) {
	return service.savedPostRepository.FindAllCollectionNames(loggedInUserID.String())
}

func (service *SavedPostService) GetAllByLoggedInUser(loggedInUserID uuid.UUID) ([]payload.SavedPostView, error) {
	savedPosts, err := service.savedPostRepository.FindAllByUserID(loggedInUserID.String())

	if err != nil {
		return nil, err
	}

	var savedPostsView []payload.SavedPostView

	for _, post := range savedPosts {
		postView := &payload.PostView{
			ID:               post.PostID,
			UserID:           post.UserID,
			Username:         service.getUserDetails(post.UserID).Username,
			Content:          post.Post.Content,
			NumberOfLikes:    service.reviewRepository.FindCountByPostIDAndStatus(post.PostID.String(), model.LIKE),
			NumberOfDislikes: service.reviewRepository.FindCountByPostIDAndStatus(post.PostID.String(), model.DISLIKE),
			NumberOfComments: service.commentRepository.FindCountByPostID(post.PostID.String()),
			Status:           service.reviewRepository.FindStatusByPostIDAndUserID(post.PostID.String(), post.UserID.String()),
			Location:         post.Post.Location,
			Description:      post.Post.Description,
		}
		savedPostView := &payload.SavedPostView{
			PostView:       *postView,
			CollectionName: post.CollectionName,
		}

		savedPostsView = append(savedPostsView, *savedPostView)
	}

	return savedPostsView, nil
}

func (service *SavedPostService) getUserDetails(userID uuid.UUID) *payload.UserDetails {
	var userIDs = &payload.UserIDs{}
	userIDs.IDs = append(userIDs.IDs, payload.UserID{ID: userID})

	requestJSON, _ := json.Marshal(userIDs)

	requestURL := fmt.Sprintf("http://%s:%s/users-details", os.Getenv("USER_SERVICE_DOMAIN"), os.Getenv("USER_SERVICE_PORT"))
	response, err := http.Post(requestURL, "application/json", bytes.NewBuffer(requestJSON))

	if err != nil {
		return nil
	}

	var details = &payload.UsersDetails{}

	helpers.FromJSON(details, response.Body)

	return &details.UsersDetails[0]
}
