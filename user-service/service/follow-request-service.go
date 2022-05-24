package service

import (
	"github.com/KristijanPill/Nishtagram/user-service/model"
	"github.com/KristijanPill/Nishtagram/user-service/repository"
	"github.com/google/uuid"
)

type FollowRequestService struct {
	requestRepository *repository.FollowRequestRepository
	followRepository  *repository.FollowRepository
}

func NewFollowRequestService(requestRepository *repository.FollowRequestRepository, followRepository *repository.FollowRepository) *FollowRequestService {
	return &FollowRequestService{
		requestRepository: requestRepository,
		followRepository:  followRepository,
	}
}

func (service *FollowRequestService) Accept(userID uuid.UUID, followerID uuid.UUID) error {
	follow := &model.Follow{
		UserID:     userID,
		FollowerID: followerID,
	}

	_, err := service.followRepository.Create(follow)

	if err != nil {
		return err
	}

	request := &model.FollowRequest{
		UserID:     followerID,
		FollowedID: userID,
	}

	return service.requestRepository.Delete(request)
}

func (service *FollowRequestService) Decline(userID uuid.UUID, followerID uuid.UUID) error {
	request := &model.FollowRequest{
		UserID:     followerID,
		FollowedID: userID,
	}

	return service.requestRepository.Delete(request)
}
