package service

import (
	"errors"

	"github.com/KristijanPill/Nishtagram/user-service/model"
	"github.com/KristijanPill/Nishtagram/user-service/payload"
	"github.com/KristijanPill/Nishtagram/user-service/repository"
	"github.com/google/uuid"
)

type FollowService struct {
	followRepository        *repository.FollowRepository
	followRequestRepository *repository.FollowRequestRepository
	userRepository          *repository.UserRepository
}

func NewFollowService(followRepository *repository.FollowRepository,
	followRequestRepository *repository.FollowRequestRepository,
	userRepository *repository.UserRepository) *FollowService {
	return &FollowService{
		followRepository:        followRepository,
		followRequestRepository: followRequestRepository,
		userRepository:          userRepository,
	}
}

func (service *FollowService) Follow(userID uuid.UUID, followedUserID uuid.UUID) (*model.Follow, error) {
	if userID == followedUserID {
		return nil, errors.New("cannot follow self")
	}

	user, err := service.userRepository.FindByID(followedUserID.String())

	if err != nil {
		return nil, err
	}

	if user.Private {
		request := &model.FollowRequest{
			UserID:     userID,
			FollowedID: followedUserID,
		}
		_, err := service.followRequestRepository.Create(request)

		if err != nil {
			return nil, err
		} else {
			return nil, nil
		}
	}

	follow := &model.Follow{
		UserID:      followedUserID,
		FollowerID:  userID,
		CloseFriend: false,
		Muted:       false,
	}

	return service.followRepository.Create(follow)
}

func (service *FollowService) UpdateCloseFriend(closeFriendID uuid.UUID, userID uuid.UUID, closeFriendStatus bool) error {
	if !service.followRepository.ExistsByUserIDAndFollowerID(userID.String(), closeFriendID.String()) {
		return errors.New("user doesnt follow you")
	}

	follow, err := service.followRepository.FindByUserIDAndFollowerID(userID.String(), closeFriendID.String())

	if err != nil {
		return err
	}

	if follow.CloseFriend == closeFriendStatus {
		return errors.New("nothing to update")
	}

	follow.CloseFriend = closeFriendStatus

	_, err = service.followRepository.Update(follow)

	return err
}

func (service *FollowService) FindByUserIDAndFollowerID(userID uuid.UUID, followerID uuid.UUID) (*model.Follow, error) {
	return service.followRepository.FindByUserIDAndFollowerID(userID.String(), followerID.String())
}

func (service *FollowService) BindFollowStatus(usersDetails *payload.UsersDetails, loggedInUserID uuid.UUID) *payload.UsersDetails {
	var retVal = []payload.UserDetails{}
	for _, details := range usersDetails.UsersDetails {
		retVal = append(retVal, payload.UserDetails{
			ID:             details.ID,
			Username:       details.Username,
			Private:        details.Private,
			ProfilePicture: details.ProfilePicture,
			Followed:       true,
		})
	}

	return &payload.UsersDetails{UsersDetails: retVal}
}

func (service *FollowService) UpdateMuted(userID uuid.UUID, followerID uuid.UUID, muted bool) error {
	if !service.followRepository.ExistsByUserIDAndFollowerID(userID.String(), followerID.String()) {
		return errors.New("user not followed")
	}

	follow, err := service.followRepository.FindByUserIDAndFollowerID(userID.String(), followerID.String())

	if err != nil {
		return err
	}

	if follow.Muted == muted {
		return errors.New("nothing to update")
	}

	follow.Muted = muted

	_, err = service.followRepository.Update(follow)

	return err
}
