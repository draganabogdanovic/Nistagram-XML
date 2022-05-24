package repository

import (
	"github.com/KristijanPill/Nishtagram/auth-service/payload"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	database *gorm.DB
}

func NewRefreshTokenRepository(database *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{database: database}
}

func (repository *RefreshTokenRepository) Create(refreshToken *payload.RefreshToken) error {
	result := repository.database.Create(refreshToken)

	return result.Error
}

func (repository *RefreshTokenRepository) Update(updatedToken *payload.RefreshToken) (*payload.RefreshToken, error) {
	result := repository.database.Save(updatedToken)

	return updatedToken, result.Error
}

func (repository *RefreshTokenRepository) FindByUserID(userID string) (payload.RefreshToken, error) {
	var refreshToken payload.RefreshToken
	result := repository.database.First(&refreshToken, "user_id = ?", userID)

	return refreshToken, result.Error
}

func (repository *RefreshTokenRepository) ExistsByUserID(userID string) bool {
	var refreshToken payload.RefreshToken
	return repository.database.Where("user_id = ?", userID).First(&refreshToken).RowsAffected == 1
}
