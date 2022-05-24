package handler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/KristijanPill/Nishtagram/story-service/helpers"
	"github.com/KristijanPill/Nishtagram/story-service/middleware"
	"github.com/KristijanPill/Nishtagram/story-service/model"
	"github.com/KristijanPill/Nishtagram/story-service/payload"
	"github.com/KristijanPill/Nishtagram/story-service/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type StoryHandler struct {
	service *service.StoryService
}

func NewStoryHandler(service *service.StoryService) *StoryHandler {
	return &StoryHandler{service: service}
}

func (handler *StoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	r.Body = ioutil.NopCloser(bytes.NewReader(body))
	requestURL := fmt.Sprintf("http://%s:%s/upload/story", os.Getenv("MEDIA_SERVICE_DOMAIN"), os.Getenv("MEDIA_SERVICE_PORT"))
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

	storyPaths := &payload.StoryUploadResponse{}

	err = helpers.FromJSON(&storyPaths, response.Body)
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

	story := &model.Story{
		UserID:           userID,
		Content:          storyPaths.StoryPaths,
		CloseFriendsOnly: r.FormValue("close_friends") == "1",
	}

	_, err = handler.service.Create(story)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func (handler *StoryHandler) FindByUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDString := vars["id"]
	userID, err := uuid.Parse(userIDString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	var loggedInUserID uuid.UUID
	if r.Context().Value(middleware.LoggedInUser{}) != nil {
		loggedInUserID = r.Context().Value(middleware.LoggedInUser{}).(uuid.UUID)
	}

	var tokenString string
	if loggedInUserID != uuid.Nil {
		tokenString = helpers.ExtractTokenFromHeader(r.Header["Authorization"][0])
	}

	stories, err := handler.service.FindByUser(userID, loggedInUserID, tokenString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	helpers.ToJSON(&stories, w)
}

func (handler *StoryHandler) FindByLoggedInUser(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.TokenKey{}).(jwt.MapClaims)

	userIDString := helpers.ExtractClaim("sub", claims)
	userID, err := uuid.Parse(userIDString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	stories, err := handler.service.FindByLoggedInUser(userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	helpers.ToJSON(&stories, w)
}

func (handler *StoryHandler) FindAllByLoggedInUser(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.TokenKey{}).(jwt.MapClaims)

	userIDString := helpers.ExtractClaim("sub", claims)
	userID, err := uuid.Parse(userIDString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	stories, err := handler.service.FindAllByLoggedInUser(userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	helpers.ToJSON(&stories, w)
}

func (handler *StoryHandler) CreateReport(w http.ResponseWriter, r *http.Request) {

	dto := &payload.ReportCreate{}
	err := helpers.FromJSON(&dto, r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	claims := r.Context().Value(middleware.TokenKey{}).(jwt.MapClaims)

	userIDString := helpers.ExtractClaim("sub", claims)
	userID, err := uuid.Parse(userIDString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	dto.UserID = userID

	_, err = handler.service.CreateReport(dto)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

}
