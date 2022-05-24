package repository

import (
	"github.com/KristijanPill/Nishtagram/user-service/model"
	"gorm.io/gorm"
)

type VerificationRequestRepository struct {
	database *gorm.DB
}

func NewVerificationRequestRepository(database *gorm.DB) *VerificationRequestRepository {
	return &VerificationRequestRepository{database: database}
}

func (repository *VerificationRequestRepository) Create(verificationRequest *model.VerificationRequest) (*model.VerificationRequest, error) {
	result := repository.database.Create(verificationRequest)

	return verificationRequest, result.Error
}

func (repository *VerificationRequestRepository) FindByVerified(verified bool) (*model.VerificationRequest, error) {
	var verReq model.VerificationRequest
	result := repository.database.First(&verReq, "verified = ?", verified)

	return &verReq, result.Error
}