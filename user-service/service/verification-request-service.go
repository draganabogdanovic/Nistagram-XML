package service

import (
	"github.com/KristijanPill/Nishtagram/user-service/model"
	"github.com/KristijanPill/Nishtagram/user-service/repository"
)

type VerificationRequestService struct {
	verificationRequestRepository *repository.VerificationRequestRepository
}

func NewVerificationRequestService(verificationRequestRepository *repository.VerificationRequestRepository) *VerificationRequestService {
	return &VerificationRequestService{
		verificationRequestRepository:  verificationRequestRepository,
	}
}

func (service *VerificationRequestService) Create(verificationRequest *model.VerificationRequest) (*model.VerificationRequest, error) {

	return service.verificationRequestRepository.Create(verificationRequest)
}