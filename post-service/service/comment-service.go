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

type CommentService struct {
	repository *repository.CommentRepository
}

func NewCommentService(repository *repository.CommentRepository) *CommentService {
	return &CommentService{repository: repository}
}

func (service *CommentService) Create(dto *payload.CommentCreate) (*model.Comment, error) {
	comment := &model.Comment{
		PostID:             dto.PostID,
		UserID:             dto.UserID,
		Content:            dto.Content,
		RepliedToCommentID: dto.RepliedToCommentID,
	}

	return service.repository.Create(comment)
}

func (service *CommentService) FindAllByPostID(id uuid.UUID) ([]payload.CommentView, error) {
	comments, err := service.repository.FindAllByPostID(id.String())

	if err != nil {
		return nil, err
	}

	var userIDs = &payload.UserIDs{}

	for _, comment := range comments {
		userIDs.IDs = append(userIDs.IDs, payload.UserID{ID: comment.UserID})
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

	var commentsView = []payload.CommentView{}

	for i, comment := range comments {
		commentView := payload.CommentView{
			ID:                 comment.ID,
			UserID:             comment.UserID,
			Username:           details.UsersDetails[i].Username,
			ProfilePicture:     details.UsersDetails[i].ProfilePicture,
			Content:            comment.Content,
			RepliedToCommentID: comment.RepliedToCommentID,
		}

		commentsView = append(commentsView, commentView)
	}

	return commentsView, nil
}
