package service

import (
	"errors"

	"github.com/KristijanPill/Nishtagram/user-service/model"
	"github.com/KristijanPill/Nishtagram/user-service/repository"
	"github.com/google/uuid"
)

type BlockService struct {
	blockRepository         *repository.BlockRepository
	followRepository        *repository.FollowRepository
	followRequestRepository *repository.FollowRequestRepository
}

func NewBlockService(blockRepository *repository.BlockRepository, followRepository *repository.FollowRepository, followRequestRepository *repository.FollowRequestRepository) *BlockService {
	return &BlockService{
		blockRepository:         blockRepository,
		followRepository:        followRepository,
		followRequestRepository: followRequestRepository,
	}
}

func (service *BlockService) Block(userID uuid.UUID, blockedID uuid.UUID) error {
	if userID == blockedID {
		return errors.New("cannot block self")
	}

	_, err := service.blockRepository.Create(&model.Block{
		UserID:    userID,
		BlockedID: blockedID,
	})

	if err != nil {
		return err
	}

	if service.followRepository.ExistsByUserIDAndFollowerID(blockedID.String(), userID.String()) {
		follow, err := service.followRepository.FindByUserIDAndFollowerID(blockedID.String(), userID.String())
		if err != nil {
			return err
		}

		err = service.followRepository.Delete(follow)

		if err != nil {
			return err
		}
	}

	if service.followRepository.ExistsByUserIDAndFollowerID(userID.String(), blockedID.String()) {
		follow, err := service.followRepository.FindByUserIDAndFollowerID(userID.String(), blockedID.String())
		if err != nil {
			return err
		}

		err = service.followRepository.Delete(follow)

		if err != nil {
			return err
		}
	}

	if service.followRequestRepository.ExistsByUserIDAndFollowedID(blockedID.String(), userID.String()) {
		request, err := service.followRequestRepository.FindByUserIDAndFollowedID(blockedID.String(), userID.String())
		if err != nil {
			return err
		}

		err = service.followRequestRepository.Delete(request)

		if err != nil {
			return err
		}
	}

	if service.followRequestRepository.ExistsByUserIDAndFollowedID(userID.String(), blockedID.String()) {
		request, err := service.followRequestRepository.FindByUserIDAndFollowedID(userID.String(), blockedID.String())
		if err != nil {
			return err
		}

		err = service.followRequestRepository.Delete(request)

		if err != nil {
			return err
		}
	}

	return nil
}

func (service *BlockService) Unblock(userID uuid.UUID, blockedID uuid.UUID) error {
	if service.blockRepository.ExistsByUserIDAndBlockedID(userID.String(), blockedID.String()) {
		block, err := service.blockRepository.FindByUserIDAndBlockedID(userID.String(), blockedID.String())

		if err != nil {
			return err
		}

		return service.blockRepository.Delete(block)
	} else {
		return errors.New("user not blocked")
	}
}
