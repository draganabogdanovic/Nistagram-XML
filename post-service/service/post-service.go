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

type PostService struct {
	postRepository    *repository.PostRepository
	reviewRepository  *repository.ReviewRepository
	commentRepository *repository.CommentRepository
}

func NewPostService(postRepository *repository.PostRepository, reviewRepository *repository.ReviewRepository, commentRepository *repository.CommentRepository) *PostService {
	return &PostService{postRepository: postRepository, reviewRepository: reviewRepository, commentRepository: commentRepository}
}

func (service *PostService) Create(post *model.Post) (*model.Post, error) {
	return service.postRepository.Create(post)
}

func (service *PostService) FindByOtherUser(userID uuid.UUID, loggedInUserID uuid.UUID) ([]payload.PostView, error) {
	posts, err := service.postRepository.FindByUserID(userID.String())

	if err != nil {
		return nil, err
	}

	var postsView []payload.PostView

	for _, post := range posts {
		postView := &payload.PostView{
			ID:               post.ID,
			UserID:           post.UserID,
			Username:         service.getUserDetails(post.UserID).Username,
			Content:          post.Content,
			NumberOfLikes:    service.reviewRepository.FindCountByPostIDAndStatus(post.ID.String(), model.LIKE),
			NumberOfDislikes: service.reviewRepository.FindCountByPostIDAndStatus(post.ID.String(), model.DISLIKE),
			NumberOfComments: service.commentRepository.FindCountByPostID(post.ID.String()),
			Status:           service.reviewRepository.FindStatusByPostIDAndUserID(post.ID.String(), loggedInUserID.String()),
			Location:         post.Location,
			Description:      post.Description,
		}
		postsView = append(postsView, *postView)
	}

	return postsView, nil
}

func (service *PostService) getUserDetails(userID uuid.UUID) *payload.UserDetails {
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

func (service *PostService) GetReviewsByUserIDAndStatus(userID uuid.UUID, status int) ([]payload.PostView, error) {
	reviews, err := service.reviewRepository.GetReviewsByUserIDAndStatus(userID, status)

	if err != nil {
		return nil, err
	}

	var postsView []payload.PostView

	for _, review := range reviews {
		postView := &payload.PostView{
			ID:               review.Post.ID,
			UserID:           review.Post.UserID,
			Username:         service.getUserDetails(review.Post.UserID).Username,
			Content:          review.Post.Content,
			NumberOfLikes:    service.reviewRepository.FindCountByPostIDAndStatus(review.Post.ID.String(), model.LIKE),
			NumberOfDislikes: service.reviewRepository.FindCountByPostIDAndStatus(review.Post.ID.String(), model.DISLIKE),
			NumberOfComments: service.commentRepository.FindCountByPostID(review.Post.ID.String()),
			Status:           service.reviewRepository.FindStatusByPostIDAndUserID(review.Post.ID.String(), userID.String()),
			Location:         review.Post.Location,
			Description:      review.Post.Description,
		}
		postsView = append(postsView, *postView)
	}

	return postsView, nil

}

func (service *PostService) CreateReport(dto *payload.ReportCreate) (*model.Report, error) {

	report := &model.Report{
		PostID:      dto.PostID,
		UserID:      dto.UserID,
		Description: dto.Description,
	}

	return service.postRepository.CreateReport(report)
}

func (service *PostService) SearchPostsByLocation(query string, loggedInUserID uuid.UUID) ([]payload.PostView, error) {
	posts, err := service.postRepository.FindLikeLocation(query)

	if err != nil {
		return nil, err
	}

	var userIDs = &payload.UserIDs{}

	for _, post := range posts {
		userIDs.IDs = append(userIDs.IDs, payload.UserID{ID: post.UserID})
	}

	requestJSON, err := json.Marshal(userIDs)

	if err != nil {
		return nil, err
	}

	requestURL := fmt.Sprintf("http://%s:%s/users-details", os.Getenv("USER_SERVICE_DOMAIN"), os.Getenv("USER_SERVICE_PORT"))
	response, err := http.Post(requestURL, "application/json", bytes.NewBuffer(requestJSON))

	if err != nil {
		return nil, err
	}

	var details = &payload.UsersDetails{}

	helpers.FromJSON(details, response.Body)

	var postsView []payload.PostView

	for i, post := range posts {
		if !details.UsersDetails[i].Private {
			postView := &payload.PostView{
				ID:               post.ID,
				UserID:           post.UserID,
				Username:         service.getUserDetails(post.UserID).Username,
				ProfilePicture:   service.getUserDetails(post.UserID).ProfilePicture,
				Content:          post.Content,
				NumberOfLikes:    service.reviewRepository.FindCountByPostIDAndStatus(post.ID.String(), model.LIKE),
				NumberOfDislikes: service.reviewRepository.FindCountByPostIDAndStatus(post.ID.String(), model.DISLIKE),
				NumberOfComments: service.commentRepository.FindCountByPostID(post.ID.String()),
				Status:           service.reviewRepository.FindStatusByPostIDAndUserID(post.ID.String(), loggedInUserID.String()),
				Location:         post.Location,
				Description:      post.Description,
			}
			postsView = append(postsView, *postView)
		}
	}

	return postsView, nil
}

func (service *PostService) SearchPostsByTags(query string, loggedInUserID uuid.UUID) ([]payload.PostView, error) {
	posts, err := service.postRepository.FindLikeTags(query)

	if err != nil {
		return nil, err
	}

	var userIDs = &payload.UserIDs{}

	for _, post := range posts {
		userIDs.IDs = append(userIDs.IDs, payload.UserID{ID: post.UserID})
	}

	requestJSON, err := json.Marshal(userIDs)

	if err != nil {
		return nil, err
	}

	requestURL := fmt.Sprintf("http://%s:%s/users-details", os.Getenv("USER_SERVICE_DOMAIN"), os.Getenv("USER_SERVICE_PORT"))
	response, err := http.Post(requestURL, "application/json", bytes.NewBuffer(requestJSON))

	if err != nil {
		return nil, err
	}

	var details = &payload.UsersDetails{}

	helpers.FromJSON(details, response.Body)

	var postsView []payload.PostView

	for i, post := range posts {
		if !details.UsersDetails[i].Private {
			postView := &payload.PostView{
				ID:               post.ID,
				UserID:           post.UserID,
				Username:         service.getUserDetails(post.UserID).Username,
				ProfilePicture:   service.getUserDetails(post.UserID).ProfilePicture,
				Content:          post.Content,
				NumberOfLikes:    service.reviewRepository.FindCountByPostIDAndStatus(post.ID.String(), model.LIKE),
				NumberOfDislikes: service.reviewRepository.FindCountByPostIDAndStatus(post.ID.String(), model.DISLIKE),
				NumberOfComments: service.commentRepository.FindCountByPostID(post.ID.String()),
				Status:           service.reviewRepository.FindStatusByPostIDAndUserID(post.ID.String(), loggedInUserID.String()),
				Location:         post.Location,
				Description:      post.Description,
			}
			postsView = append(postsView, *postView)
		}
	}

	return postsView, nil
}
