package repository

import (
	"github.com/KristijanPill/Nishtagram/user-service/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	database *gorm.DB
}

func NewUserRepository(database *gorm.DB) *UserRepository {
	return &UserRepository{database: database}
}

func (repository *UserRepository) Create(user *model.User) (*model.User, error) {
	result := repository.database.Create(user)

	return user, result.Error
}

func (repository *UserRepository) FindByID(id string) (*model.User, error) {
	var user model.User
	result := repository.database.First(&user, "id = ?", id)

	return &user, result.Error
}

func (repository *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	result := repository.database.First(&user, "username = ?", username)

	return &user, result.Error
}

func (repository *UserRepository) Update(updatedUser *model.User) (*model.User, error) {
	result := repository.database.Save(updatedUser)

	return updatedUser, result.Error
}

func (repository *UserRepository) IsUsernameUnique(username string) bool {
	var user model.User
	if err := repository.database.First(&user, "username = ?", username).Error; err != nil {
		return true
	}

	return false
}

func (repository *UserRepository) FindLikeQuery(query string) ([]model.User, error) {
	var users []model.User
	result := repository.database.Where("username ILIKE ? AND private = ?", "%"+query+"%", false).Or("name ILIKE ? AND private = ?", "%"+query+"%", false).Find(&users)

	return users, result.Error
}

func (repository *UserRepository) Delete(id string) error {
	result := repository.database.Delete(&model.User{}, id)

	return result.Error
}
