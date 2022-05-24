package service

import (
	"github.com/KristijanPill/Nishtagram/post-service/model"
	"github.com/KristijanPill/Nishtagram/post-service/payload"
	"github.com/KristijanPill/Nishtagram/post-service/repository"
)

type ReviewService struct {
	repository *repository.ReviewRepository
}

func NewReviewService(repository *repository.ReviewRepository) *ReviewService {
	return &ReviewService{repository: repository}
}

func (service *ReviewService) ReviewPost(dto *payload.ReviewCreate) (*model.Review, error) {
	exists, err := service.repository.ExistsByPostIDAndUserID(dto.PostID.String(), dto.UserID.String())

	if err != nil {
		return nil, err
	}

	review := &model.Review{
		PostID: dto.PostID,
		UserID: dto.UserID,
		Status: dto.ReviewStatus,
	}

	if exists {
		if review.Status == model.NONE {

			return nil, service.repository.Delete(review.PostID.String(), review.UserID.String())
		} else {
			return service.repository.Update(review)
		}
	} else {
		return service.repository.Create(review)
	}
}
