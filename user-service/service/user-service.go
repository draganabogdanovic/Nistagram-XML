package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/KristijanPill/Nishtagram/user-service/model"
	"github.com/KristijanPill/Nishtagram/user-service/payload"
	"github.com/KristijanPill/Nishtagram/user-service/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository   *repository.UserRepository
	followRepository *repository.FollowRepository
}

func NewUserService(userRepository *repository.UserRepository, followRepository *repository.FollowRepository) *UserService {
	return &UserService{
		userRepository:   userRepository,
		followRepository: followRepository,
	}
}

func (service *UserService) Create(dto *payload.CreateUser) (*model.User, error) {
	if !service.userRepository.IsUsernameUnique(dto.Username) {
		return nil, errors.New("Username unavailable.")
	}

	user := &model.User{
		Username:               dto.Username,
		Role:                   model.USER,
		Email:                  dto.Email,
		Private:                false,
		Verified:               false,
		Taggable:               true,
		CanRecieveAnonMessages: true,
		Name:                   dto.Name,
		DOB:                    dto.DOB,
		Gender:                 dto.Gender,
		PhoneNumber:            dto.PhoneNumber,
		Website:                dto.Website,
		Bio:                    dto.Bio,
		CreatedAt:              time.Now(),
	}

	user, err := service.userRepository.Create(user)

	if err != nil {
		return nil, err
	}

	credentials := &payload.Credentials{
		ID:       user.ID,
		Username: dto.Username,
		Password: dto.Password,
	}

	requestJSON, _ := json.Marshal(credentials)
	requestURL := fmt.Sprintf("http://%s:%s/register", os.Getenv("AUTH_SERVICE_DOMAIN"), os.Getenv("AUTH_SERVICE_PORT"))

	response, err := http.Post(requestURL, "application/json", bytes.NewBuffer(requestJSON))

	if err != nil {
		service.userRepository.Delete(user.ID.String())
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		service.userRepository.Delete(user.ID.String())

		return nil, errors.New("Could not register user.")
	}

	return user, nil
}

func (service *UserService) FindByID(id uuid.UUID) (*model.User, error) {
	return service.userRepository.FindByID(id.String())
}

func (service *UserService) FindByUsername(username string) (*model.User, error) {
	return service.userRepository.FindByUsername(username)
}

func (service *UserService) Update(dto *payload.UserInfo, id uuid.UUID) (*model.User, error) {
	user, err := service.userRepository.FindByID(id.String())

	if err != nil {
		return nil, err
	}

	user.Email = dto.Email
	user.Name = dto.Name
	user.DOB = dto.DOB
	user.Gender = dto.Gender
	user.PhoneNumber = dto.PhoneNumber
	user.Website = dto.Website
	user.Bio = dto.Bio

	return service.userRepository.Update(user)
}

func (service *UserService) BindUsernameToID(userIDs *payload.UserIDs) *payload.UsersDetails {
	var usersDetails = []payload.UserDetails{}

	for _, userID := range userIDs.IDs {
		user, err := service.userRepository.FindByID(userID.ID.String())
		details := payload.UserDetails{
			ID:             user.ID,
			Username:       user.Username,
			Private:        user.Private,
			ProfilePicture: user.ProfilePicture,
		}

		if err != nil {
			usersDetails = append(usersDetails, payload.UserDetails{})
			continue
		}

		usersDetails = append(usersDetails, details)
	}

	return &payload.UsersDetails{UsersDetails: usersDetails}
}

func (service *UserService) UpdateProfilePicture(profilePicturePath string, userID uuid.UUID) (*model.User, error) {
	user, err := service.userRepository.FindByID(userID.String())

	if err != nil {
		return nil, err
	}

	user.ProfilePicture = profilePicturePath

	service.userRepository.Update(user)

	return user, nil
}

func (service *UserService) hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 14)
}

func (service *UserService) GetFollowerCount(id uuid.UUID) int64 {
	return service.followRepository.GetFollowerCount(id.String())
}

func (service *UserService) GetFollowingCount(id uuid.UUID) int64 {
	return service.followRepository.GetFollowingCount(id.String())
}

func (service *UserService) SearchUsers(query string) ([]payload.UserView, error) {
	result, err := service.userRepository.FindLikeQuery(query)

	if err != nil {
		return nil, err
	}

	var users = []payload.UserView{}

	for _, user := range result {
		users = append(users, payload.UserView{
			Username:       user.Username,
			Name:           user.Name,
			ProfilePicture: user.ProfilePicture,
		})
	}

	return users, nil
}
