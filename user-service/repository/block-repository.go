package repository

import (
	"github.com/KristijanPill/Nishtagram/user-service/model"
	"gorm.io/gorm"
)

type BlockRepository struct {
	database *gorm.DB
}

func NewBlockRepository(database *gorm.DB) *BlockRepository {
	return &BlockRepository{database: database}
}

func (repository *BlockRepository) Create(block *model.Block) (*model.Block, error) {
	result := repository.database.Create(block)

	return block, result.Error
}

func (repository *BlockRepository) ExistsByUserIDAndBlockedID(userID string, blockedID string) bool {
	var block model.Block
	return repository.database.Where("user_id = ? AND blocked_id = ?", userID, blockedID).First(&block).RowsAffected == 1
}

func (repository *BlockRepository) FindByUserIDAndBlockedID(userID string, blockedID string) (*model.Block, error) {
	var block model.Block
	result := repository.database.Where("user_id = ? AND blocked_id = ?", userID, blockedID).First(&block)

	return &block, result.Error
}

func (repository *BlockRepository) Delete(block *model.Block) error {
	result := repository.database.Delete(block)

	return result.Error
}
