package repository

import (
	"github.com/KristijanPill/Nishtagram/auth-service/payload"
	"gorm.io/gorm"
)

type CredentialsRepository struct {
	database *gorm.DB
}

func NewCredentialsRepository(database *gorm.DB) *CredentialsRepository {
	return &CredentialsRepository{database: database}
}

func (repository *CredentialsRepository) Create(credentials *payload.Credentials) error {
	result := repository.database.Create(credentials)

	return result.Error
}

func (repository *CredentialsRepository) FindByUsername(username string) (payload.Credentials, error) {
	var credentials payload.Credentials
	result := repository.database.First(&credentials, "username = ?", username)

	return credentials, result.Error
}
