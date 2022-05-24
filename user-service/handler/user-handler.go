package handler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/KristijanPill/Nishtagram/user-service/helpers"
	"github.com/KristijanPill/Nishtagram/user-service/middleware"
	"github.com/KristijanPill/Nishtagram/user-service/payload"
	"github.com/KristijanPill/Nishtagram/user-service/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	userService   *service.UserService
	followService *service.FollowService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{userService: service}
}

func (handler *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	dto := &payload.CreateUser{}
	err := helpers.FromJSON(&dto, r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	_, err = handler.userService.Create(dto)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
}

func (handler *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.TokenKey{}).(jwt.MapClaims)

	userIDString := helpers.ExtractClaim("sub", claims)
	userID, err := uuid.Parse(userIDString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	dto := &payload.UserInfo{}
	err = helpers.FromJSON(&dto, r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	user, err := handler.userService.Update(dto, userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	fmt.Println(user)
}

func (handler *UserHandler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.TokenKey{}).(jwt.MapClaims)

	userIDString := helpers.ExtractClaim("sub", claims)
	userID, err := uuid.Parse(userIDString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	user, err := handler.userService.FindByID(userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	userInfo := &payload.UserInfo{
		Username:       user.Username,
		Email:          user.Email,
		Name:           user.Name,
		DOB:            user.DOB,
		Gender:         user.Gender,
		PhoneNumber:    user.PhoneNumber,
		Website:        user.Website,
		Bio:            user.Bio,
		ProfilePicture: user.ProfilePicture,
	}

	helpers.ToJSON(&userInfo, w)
}

func (handler *UserHandler) GetUserProfileInfo(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.TokenKey{}).(jwt.MapClaims)

	userIDString := helpers.ExtractClaim("sub", claims)
	userID, err := uuid.Parse(userIDString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	user, err := handler.userService.FindByID(userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	userProfileInfo := &payload.UserProfileInfo{
		Username:       user.Username,
		Name:           user.Name,
		Bio:            user.Bio,
		Website:        user.Website,
		Followers:      handler.userService.GetFollowerCount(user.ID),
		Following:      handler.userService.GetFollowingCount(user.ID),
		ProfilePicture: user.ProfilePicture,
	}

	helpers.ToJSON(&userProfileInfo, w)
}

func (handler *UserHandler) GetOtherUserProfileInfo(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	username := vars["username"]

	user, err := handler.userService.FindByUsername(username)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	userProfileInfo := &payload.UserProfileInfo{
		ID:             user.ID,
		Username:       user.Username,
		Name:           user.Name,
		Bio:            user.Bio,
		Website:        user.Website,
		Followers:      handler.userService.GetFollowerCount(user.ID),
		Following:      handler.userService.GetFollowingCount(user.ID),
		ProfilePicture: user.ProfilePicture,
	}

	helpers.ToJSON(&userProfileInfo, w)
}

func (handler *UserHandler) BindUsernameToID(w http.ResponseWriter, r *http.Request) {
	var userIDs = &payload.UserIDs{}
	err := helpers.FromJSON(userIDs, r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	usernames := handler.userService.BindUsernameToID(userIDs)

	helpers.ToJSON(usernames, w)
}

func (handler *UserHandler) UpdateProfilePicture(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	r.Body = ioutil.NopCloser(bytes.NewReader(body))
	requestURL := fmt.Sprintf("http://%s:%s/upload/profile-picture", os.Getenv("MEDIA_SERVICE_DOMAIN"), os.Getenv("MEDIA_SERVICE_PORT"))
	proxyReq, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(body))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	for header, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(header, value)
		}
	}

	client := &http.Client{}
	response, err := client.Do(proxyReq)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	profilePicturePath := &payload.ProfilePictureUploadResponse{}

	err = helpers.FromJSON(&profilePicturePath, response.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	claims := r.Context().Value(middleware.TokenKey{}).(jwt.MapClaims)

	userIDString := helpers.ExtractClaim("sub", claims)
	userID, err := uuid.Parse(userIDString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	_, err = handler.userService.UpdateProfilePicture(profilePicturePath.ProfilePicture, userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func (handler *UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")

	result, err := handler.userService.SearchUsers(query)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	helpers.ToJSON(&result, w)
}

func (handler *UserHandler) GetUsersDetails(w http.ResponseWriter, r *http.Request) {
	var loggedInUserID uuid.UUID

	if r.Context().Value(middleware.LoggedInUser{}) != nil {
		loggedInUserID = r.Context().Value(middleware.LoggedInUser{}).(uuid.UUID)
	}

	var userIDs = &payload.UserIDs{}
	err := helpers.FromJSON(userIDs, r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	usersDetails := handler.userService.BindUsernameToID(userIDs)
	details := handler.followService.BindFollowStatus(usersDetails, loggedInUserID)

	helpers.ToJSON(details, w)
}
